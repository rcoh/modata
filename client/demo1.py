import coding
import restlib

data_chunks = [str(x) for x in xrange(10)]
chunks_meta = coding.make_chunk_dicts(data_chunks)

for chunk in chunks_meta:
    restlib.store(chunk['data'], chunk['digest'])

print len(chunks_meta), "items stored"
success = True
for chunk in chunks_meta:
    data = restlib.findvalue(chunk['digest'])
    if data != chunk['data']:
        success = False
        print "DATA MISMATCH FOR DIGEST", chunk['digest']
    
if success:
    print len(chunks_meta), "items retrieved sucessfully"