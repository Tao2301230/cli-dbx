#!/usr/bin/env python3
"""Prepare synthetic, isolated fixtures for the eight-fix manual regression."""
import json
from pathlib import Path
import sqlite3
import sys

if len(sys.argv) != 2:
    raise SystemExit("usage: python3 prepare.py <new-fixture-directory>")
root = Path(sys.argv[1]).expanduser().resolve()
root.mkdir(parents=True, exist_ok=False)
profiles = root / "connections"
profiles.mkdir()
database = root / "manual.sqlite"
with sqlite3.connect(database) as db:
    db.executescript("""
        CREATE TABLE users(id INTEGER PRIMARY KEY, name TEXT, active INTEGER);
        INSERT INTO users VALUES(1,'original',1),(2,'second',1),(3,'third',0);
        CREATE TABLE forbidden(id INTEGER, token TEXT);
        INSERT INTO forbidden VALUES(1,'SYNTHETIC_PRIVATE_TOKEN');
        CREATE TABLE imports(i INTEGER, n TEXT);
    """)

def write_profile(name, body):
    (profiles / (name + ".toml")).write_text("[connection]\n" + body, encoding="utf-8")

base = 'engine="sqlite"\npath=' + json.dumps(str(database)) + '\n'
write_profile("reader", base + 'allow_actions=["query"]\nallow_tables=["users"]\n')
write_profile("writer", base + 'allow_actions=["query","update"]\nallow_tables=["users"]\n')
write_profile("all", base + 'allow_actions=["query","update","schema","admin"]\n')
write_profile("pg-reader", 'engine="postgres"\ndsn="postgres://dummy:DUMMY_SECRET@localhost/fake"\nallow_actions=["query"]\nallow_tables=["users"]\n')
write_profile("pg-writer", 'engine="postgres"\ndsn="postgres://dummy:DUMMY_SECRET@localhost/fake"\nallow_actions=["query","update"]\nallow_tables=["users"]\n')
write_profile("sqlserver", 'engine="sqlserver"\ndsn="sqlserver://dummy:OLD_DUMMY_SECRET@localhost:1433?database=fake&change+password=NEW_DUMMY_SECRET"\nallow_actions=["query"]\n')
(root / "large.json").write_text('[{"i":9007199254740993,"n":9223372036854775808}]\n')
(root / "invalid-format-plan.json").write_text(json.dumps({"steps": [{"action": "update", "sql": "update users set name='changed' where id=1"}]}))
(root / "existing.csv").write_text("KEEP_THIS_FILE\n")
print(str(root))
