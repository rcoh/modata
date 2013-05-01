import hashlib
import zfec
import math
import requests
CHUNK_SIZE = 50 
SERVER = "http://localhost:1234/"

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

def recombine_chunks(k, m, chunks):
    decer = zfec.easyfec.Dencoder(int(k), int(m))
    return decer.decode(chunks)


def make_digests(chunks):
    return [(digest_for_chunk(chunk), chunk) for chunk in chunks]

def send_chunks_to_storage(chunks):
    for (digest, chunk) in chunks:
        requests.post(SERVER + "upload", data={'key': digest, 'data': chunk})
    print "Sent", len(chunks), "chunks to datastore"
    print "Sizes were: ", [len(c) for (d, c) in chunks]

def send_chunks_to_meta(chunks):
    print "Sent", len(chunks), "chunks to meta"

def get_chunks(metadata):
    k = metadata['k']
    m = metadata['m']
    chunks = [requets.get(SERVER + "get/" + chunkhash).text for chunkhash in metadata['chunks']]
    return recombine_chunks(k, m, chunks)
