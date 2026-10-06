"""Offline negative package checks. Copies only; never alters installed drivers."""
import argparse,hashlib,json,pathlib,shutil,subprocess,tempfile
p=argparse.ArgumentParser();p.add_argument('--cli',required=True);p.add_argument('--driver',required=True);a=p.parse_args()
cli=str(pathlib.Path(a.cli).resolve());src=pathlib.Path(a.driver).resolve();manifest=json.loads(pathlib.Path(str(src)+'.json').read_text())
root=pathlib.Path(tempfile.mkdtemp(prefix='dbx-bad-package-'));results=[]
for kind,expected in [('missing','driver_missing'),('protocol','driver_protocol_mismatch'),('engine','driver_invalid'),('platform','driver_invalid'),('adapter','driver_invalid'),('checksum','driver_integrity_failed'),('resource','driver_integrity_failed')]:
 folder=root/kind;folder.mkdir();binary=folder/src.name;shutil.copy2(src,binary);m=dict(manifest);m['files']={}
 if kind=='protocol':m['protocol_version']=999
 if kind=='engine':m['engine']='wrong-engine'
 if kind=='platform':m['arch']='wrong-arch'
 if kind=='adapter':m['adapter']='native' if manifest.get('adapter')=='odbc' else 'odbc'
 if kind=='checksum':m['sha256']='0'*64
 if kind=='resource':
  (folder/'probe.txt').write_text('changed');m['files']={'probe.txt':'0'*64}
 if kind!='missing':pathlib.Path(str(binary)+'.json').write_text(json.dumps(m))
 config='[connection]\nengine='+json.dumps(manifest['engine'])+'\ndriver_path='+json.dumps(str(binary))+'\nhost="127.0.0.1"\nport=1\nuser="DUMMY"\ndatabase="DUMMY"\nallow_actions=["query"]\ntimeout="1s"\n'
 if manifest.get('adapter')=='odbc':config+='driver="odbc"\nodbc_driver='+json.dumps(str(src))+'\n'
 (folder/'probe.toml').write_text(config)
 command=[cli,'conn','test','probe','--config',str(folder/'probe.toml')]
 r=subprocess.run(command,capture_output=True,text=True,timeout=10)
 try: payload=json.loads(r.stdout)
 except ValueError:payload={}
 matched=r.returncode==1 and payload.get('code')==expected
 results.append(dict(case=kind,command=command,expected_code=expected,exit=r.returncode,matched=matched,stdout=r.stdout,stderr=r.stderr))
 print(kind, 'matched' if matched else 'MISMATCH',r.returncode,payload.get('code'))
(root/'results.json').write_text(json.dumps(results,ensure_ascii=False,indent=2))
print('Results:',root/'results.json')
raise SystemExit(0 if all(x['matched'] for x in results) else 1)
