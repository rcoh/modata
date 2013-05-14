import hashlib
import zfec
import math
import restlib
import threading
import sys
from copy import deepcopy
from multiprocessing import Pool

CHUNK_SIZE = 4096
CHUNK_SIZE = 32768
# Encoder limit
MAX_CHUNKS = 256
#TODO: zlib

def basic_chunk(data, size=CHUNK_SIZE):
    for i in xrange(0, len(data), size):
        yield data[i:i+size]

def digest_for_chunk(chunk):
    return hashlib.sha1(chunk).hexdigest()

def send_chunks_get_metadata(data):
    k,m,superblock_encodings = erasure_chunk(data)
    hexed_supers = [hexify_superblock(s) for s in superblock_encodings]
    super_dicts_chunks = []
    metadata = []
    for superb in hexed_supers:
        super_dicts_chunks.append(make_chunk_dicts(superb))

    print "Sending chunks to storage"
    for superb_dict in super_dicts_chunks:
        send_chunks_to_storage(superb_dict)
        metadata.append(get_metadata(k, m, superb_dict))

    return metadata

def hexify_superblock(block):
    return [c.encode("hex") for c in block]

def erasure_chunk(data, percent_required=.5, target_size = CHUNK_SIZE):
    # k / m == percent required
    # Variable size chunks based on target
    # 
    # k = total / target_size == pre_replication_chunks
    # total_chunks = k * repfactor = k * 2
    # 
    # Need 256 total chunks => need k = 256 / repfactor]
    # len(superblock) == (256 * factor) * chunk_size

    print len(data)
    print CHUNK_SIZE
    superblock_length = min(int(math.ceil(MAX_CHUNKS * percent_required * CHUNK_SIZE)),
                            len(data))

    trial_k = math.ceil(superblock_length / float(target_size))
    trial_m = trial_k / percent_required
    print "m for subblock: ", trial_m
    assert trial_m <= 256
    # m must be less than 256 for the encoder
    encer = zfec.easyfec.Encoder(int(trial_k), int(trial_m))
    encodings = []
    print "Total superblocks: ", len(data) / superblock_length
    for superblock in basic_chunk(data, superblock_length):
        data = superblock
        if isinstance(data, unicode):
            data = data.encode('utf-8')
        encodings.append(encer.encode(data))
    return (trial_k, trial_m, encodings)

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
        if chunk['data']:
          used_data.append(chunk['data'].decode('hex'))
          used_blocknums.append(chunk['blocknum'])
    print len(used_data)
    print k

    return ''.join(decer.decode(used_data[:k], used_blocknums[:k]))

def make_chunk_dicts(chunks):
    return [{'digest': digest_for_chunk(data), 'data': data, 'blocknum': i} for i, data in enumerate(chunks)]

def send_chunks_to_storage(chunks):
    threads = []
    for chunk_dict in chunks:
        digest = chunk_dict['digest']
        chunk = chunk_dict['data']
        t = threading.Thread(target = restlib.store,
                             args = (deepcopy(chunk), deepcopy(digest), False))
        t.daemon = True
        threads.append(t)
        t.start()
    count = len(threads)
    for t in threads:
      t.join()
      print "{} left".format(count)
      count -= 1

def get_metadata(k, m, chunk_dicts):
    meta = {}
    meta['k'] = k
    meta['m'] = m
    for chunk in chunk_dicts:
        del chunk['data']
    meta['chunks'] = chunk_dicts
    return meta

def get_data(digest):
    try:
        return restlib.findvalue(digest, local=False)
    except Exception as e:
        print e
        return None

def data_for_superblock(metadata):
    k = int(metadata['k'])
    m = int(metadata['m'])
    # TODO: we only need k
    print "Looking for {} chunks".format(len(metadata['chunks']))
    required = len(metadata['chunks'])

    pool = Pool(processes=4)
    data = pool.map(get_data,
                    [cd['digest'] for cd in metadata['chunks']])
    index = 0
    for chunk_dict in metadata['chunks']:
        chunk_dict['data'] = data[index]
        index += 1

    return recombine_chunks(k, m, metadata['chunks'])

def get_chunks(metadata):
    if type(metadata) != type([]):
        print "Bad metadata. Try again"
        return

    print "Need {} blocks".format(len(metadata))
    result = '' 
    for superblock in metadata:
        result += data_for_superblock(superblock)

    return result
