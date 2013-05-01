import hashlib
import zfec
import math
import requests
CHUNK_SIZE = 256 
SERVER = "http://localhost:1234/"

#TODO: zlib

def basic_chunk(data, size=CHUNK_SIZE):
    for i in xrange(0, len(data), size):
        yield data[i:i+size]

def digest_for_chunk(chunk):
    return hashlib.sha1(chunk).hexdigest()

def send_chunks_get_metadata(data):
    k,m,chunks = erasure_chunk(data)
    dict_chunks = make_chunk_dicts(chunks)
    send_chunks_to_storage(dict_chunks)
    return get_metadata(k, m, dict_chunks)

def erasure_chunk(data, percent_required=.5, target_size = CHUNK_SIZE):
    # k / m == percent required
    k = math.ceil(len(data) / float(target_size))
    m = math.ceil(k / percent_required)
    encer = zfec.easyfec.Encoder(int(k), int(m))
    return (k, m, encer.encode(data))

def recombine_chunks(k, m, chunks):
    k = int(k)
    m = int(m)
    decer = zfec.Decoder(k, m)
    # Pick k chunks:
    used_data = []
    used_blocknums = []
    if len(chunks) < k:
        print "Not enough chunks to recover data"
        return ""
    for chunk in chunks:
        used_data.append(chunk['data'])
        used_blocknums.append(chunk['blocknum'])

    return ''.join(decer.decode(used_data[:k], used_blocknums[:k]))

def make_chunk_dicts(chunks):
    return [{'digest': digest_for_chunk(data), 'data': data, 'blocknum': i} for i, data in enumerate(chunks)]

def send_chunks_to_storage(chunks):
    for chunk_dict in chunks:
        digest = chunk_dict['digest']
        chunk = chunk_dict['data']
        resp_obj = requests.post(SERVER + "store", data={'key': digest, 'data': chunk})
        resp = resp_obj.json
        if resp['status'] != 'OK':
            print resp

def get_metadata(k, m, chunk_dicts):
    meta = {}
    meta['k'] = k
    meta['m'] = m
    for chunk in chunk_dicts:
        del chunk['data']
    meta['chunks'] = chunk_dicts
    return meta

def get_chunks(metadata):
    k = metadata['k']
    m = metadata['m']
    # TODO: we only need k
    for chunk_dict in metadata['chunks']:
        resp = requests.get(SERVER + "find-value/" + chunk_dict['digest']).json
        if resp['status'] == 'OK':
            chunk_dict['data'] = resp['data']

    return recombine_chunks(k, m, metadata['chunks'])
