import hashlib
import zfec
import math
import requests
CHUNK_SIZE = 50 
SERVER = "http://localhost:1234/upload"

def basic_chunk(data, size=CHUNK_SIZE):
    for i in xrange(0, len(data), size):
        yield data[i:i+size]

def digest_for_chunk(chunk):
    return hashlib.sha1(chunk).hexdigest()

def erasure_chunk(data, percent_required=.5, target_size = CHUNK_SIZE):
    # k / m == percent required
    k = math.ceil(len(data) / float(target_size))
    m = math.ceil(k / percent_required)
    encer = zfec.easyfec.Encoder(int(k), int(m))
    print "k:", k, "m:", m
    return encer.encode(data)

def send_chunks_to_storage(chunks):
    requests.post(SERVER, data={'key': 'test', 'file': chunks[0]})
    print "Sent", len(chunks), "chunks to datastore"
    print "Sizes were: ", [len(c) for c in chunks]

def send_chunks_to_meta(chunks):
    digests = [digest_for_chunk(chunk) for chunk in chunks]
    print "Sent", len(chunks), "chunks to meta"
