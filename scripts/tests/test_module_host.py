"""Real shared-host smoke against a local PostgreSQL protocol fixture, not a DB.

Run with DBX_MODULE_HOST=/absolute/path/to/built/dbx python3 -m unittest
discover -s scripts/tests -p test_module_host.py -v. Modules sit beside the host.
"""

import json
import os
from pathlib import Path
import socket
import struct
import subprocess
import tempfile
import threading
import unittest


HOST = os.environ.get("DBX_MODULE_HOST")
SQL = "select 42, '12345678901234567890.1234'::numeric, 'PG动态装载'::text"


def receive(sock, length):
    chunks = bytearray()
    while len(chunks) < length:
        chunk = sock.recv(length - len(chunks))
        if not chunk:
            raise EOFError("connection closed")
        chunks.extend(chunk)
    return bytes(chunks)


def message(kind, body=b""):
    return kind + struct.pack("!I", len(body) + 4) + body


@unittest.skipUnless(HOST, "set DBX_MODULE_HOST to run the built module host")
class ModuleHostTests(unittest.TestCase):
    def test_shared_host_keeps_common_code_without_vendor_drivers(self):
        host = Path(HOST).resolve()
        go = os.environ.get("DBX_MODULE_GO", "go")
        proc = subprocess.run([go, "tool", "nm", str(host)],
                              capture_output=True, text=True, check=True)
        symbols = {line.split(maxsplit=2)[2] for line in proc.stdout.splitlines()
                   if len(line.split(maxsplit=2)) == 3}
        for name in ("crypto/tls.X509KeyPair", "crypto/x509.ParsePKCS8PrivateKey",
                     "golang.org/x/text/language.Parse",
                     "golang.org/x/text/message.(*Printer).Sprintf",
                     "math/big.(*Rat).FloatString", "encoding/json.(*Decoder).Token"):
            with self.subTest(symbol=name):
                self.assertTrue(name in symbols, f"shared host dropped {name}")
        for prefix in ("github.com/jackc/pgx/v5.", "github.com/jackc/pgx/v5/",
                       "github.com/oracle/go-oracledb/v26/", "gitee.com/chunanyong/dm.",
                       "github.com/microsoft/go-mssqldb.", "github.com/microsoft/go-mssqldb/"):
            with self.subTest(vendor=prefix):
                self.assertFalse(any(name.startswith(prefix) for name in symbols),
                                 f"shared host embedded vendor code: {prefix}")

    def test_all_modules_load(self):
        host = Path(HOST).resolve()
        for engine in ("postgres", "oracle", "dm", "sqlserver"):
            with self.subTest(engine=engine):
                proc = subprocess.run([str(host), "--internal-module-check", engine,
                                       str(host.parent / f"{engine}.dbmod")],
                                      capture_output=True, text=True, timeout=20)
                self.assertEqual(proc.returncode, 0, proc.stderr)
                result = json.loads(proc.stdout)
                self.assertTrue(result["loaded"])
                self.assertEqual(result["engine"], engine)
                self.assertFalse(result["database_connected"])

    def test_real_cli_query_and_policy(self):
        host = Path(HOST).resolve()
        module = host.parent / "postgres.dbmod"
        listener = socket.socket()
        self.addCleanup(listener.close)
        listener.bind(("127.0.0.1", 0))
        listener.listen()
        listener.settimeout(20)
        port = listener.getsockname()[1]
        errors, queries = [], []

        def serve():
            try:
                sock, _ = listener.accept()
                with sock:
                    sock.settimeout(20)
                    length = struct.unpack("!I", receive(sock, 4))[0]
                    startup = receive(sock, length - 4)
                    if struct.unpack("!I", startup[:4])[0] != 196608:
                        raise AssertionError("expected PostgreSQL v3 startup")
                    out = message(b"R", struct.pack("!I", 0))
                    for key, value in [("server_version", "16.0"), ("client_encoding", "UTF8"),
                                       ("server_encoding", "UTF8"), ("DateStyle", "ISO, MDY"),
                                       ("integer_datetimes", "on"), ("standard_conforming_strings", "on"),
                                       ("TimeZone", "UTC")]:
                        out += message(b"S", key.encode() + b"\0" + value.encode() + b"\0")
                    sock.sendall(out + message(b"K", struct.pack("!II", 1, 1)) + message(b"Z", b"I"))
                    while True:
                        try:
                            kind = receive(sock, 1)
                        except EOFError:
                            break
                        length = struct.unpack("!I", receive(sock, 4))[0]
                        body = receive(sock, length - 4)
                        if kind == b"X":
                            break
                        if kind != b"Q":
                            raise AssertionError(f"unexpected frontend message {kind!r}")
                        query = body[:-1].decode()
                        if query == "-- ping":
                            sock.sendall(message(b"I") + message(b"Z", b"I"))
                            continue
                        if query.strip().rstrip(";") != SQL:
                            raise AssertionError(f"unexpected SQL: {query!r}")
                        queries.append(query)
                        fields = b""
                        for name, oid, size in [("id", 20, 8), ("amount", 1700, -1), ("message", 25, -1)]:
                            fields += name.encode() + b"\0" + struct.pack("!IhIhih", 0, 0, oid, size, -1, 0)
                        values = [b"42", b"12345678901234567890.1234", "PG动态装载".encode()]
                        row = struct.pack("!h", len(values)) + b"".join(struct.pack("!i", len(v)) + v for v in values)
                        sock.sendall(message(b"T", struct.pack("!h", 3) + fields) + message(b"D", row)
                                     + message(b"C", b"SELECT 1\0") + message(b"Z", b"I"))
            except BaseException as exc:
                errors.append(exc)

        thread = threading.Thread(target=serve, daemon=True)
        try:
            with tempfile.TemporaryDirectory() as folder:
                config = Path(folder) / "remote.toml"
                dsn = f"postgres://app:test@127.0.0.1:{port}/app?sslmode=disable&default_query_exec_mode=simple_protocol"
                config.write_text("[connection]\nengine='postgres'\ndriver='module'\n"
                                  + f"driver_path={json.dumps(str(module))}\ndsn={json.dumps(dsn)}\n"
                                  + "allow_actions=['query']\ntimeout_seconds=10\n")
                def run(action, sql):
                    proc = subprocess.run([str(host), action, "--config", folder, "remote", sql],
                                          capture_output=True, text=True, timeout=20)
                    self.assertEqual(proc.stderr, "", proc.stderr)
                    return proc, json.loads(proc.stdout)

                blocked, payload = run("schema", "create table x (id int)")
                self.assertNotEqual(blocked.returncode, 0)
                self.assertEqual(payload["code"], "action_blocked")
                thread.start()
                proc, payload = run("query", SQL)
                self.assertEqual(proc.returncode, 0, payload)
                self.assertTrue(payload["ok"], payload)
                rows = payload["data"]["rows"]
                self.assertEqual(rows, [{"id": 42, "amount": "12345678901234567890.1234", "message": "PG动态装载"}])
                thread.join(timeout=5)
                self.assertFalse(thread.is_alive(), "module worker did not close the connection")
                self.assertEqual(errors, [])
                self.assertEqual(len(queries), 1)
        finally:
            listener.close()


if __name__ == "__main__":
    unittest.main()
