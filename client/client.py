import argparse
import lib

parser = argparse.ArgumentParser(description='Store a file in modata')
parser.add_argument('--file', help='File to upload', required=True)
args = parser.parse_args()

filename = args.file

with open(filename, 'r') as f:
    data = f.read()
    chunks = list(lib.erasure_chunk(data))
    lib.send_chunks_to_storage(chunks)
    access_info = lib.send_chunks_to_meta(chunks)
