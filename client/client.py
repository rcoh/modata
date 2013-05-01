import argparse
import coding 
import json
import keyfilelib
import sys

parser = argparse.ArgumentParser(description='modataclient')
parser.add_argument('--keyfile', dest='keyfile')
parser.add_argument('--upload', dest='upload')
parser.add_argument('--download', dest='download')

args = parser.parse_args()

if not args.keyfile:
    keyfile_name = keyfilelib.create_new_keyfile()
else:
    keyfile_name = args.keyfile

with open(keyfile_name, 'r') as keyfile_handle:
    keyfile = json.loads(keyfile_handle.read())

if args.upload:
    with open(args.upload, 'r') as f:
        data = f.read()
        metadata = coding.send_chunks_get_metadata(data)
        keyfile[args.upload] = metadata
        keyfilelib.save(keyfile_name, keyfile)

if args.download:
    if not args.download in keyfile:
        print "Unknown file. Exiting"
        sys.exit(1)
    file_metadata = keyfile[args.download]
    data = coding.get_chunks(file_metadata)
    print data

