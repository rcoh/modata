import coding
import restlib
from runner import ServerManager

manager = ServerManager()
manager.async_start_block_server(1234)
manager.async_start_block_server(1235)

try:
    data_chunks = [str(x) for x in xrange(10)]
    chunks_meta = coding.make_chunk_dicts(data_chunks)

    for chunk in chunks_meta:
        restlib.store(chunk['data'], chunk['digest'])

    print len(chunks_meta), "items stored"
    success = True
    for chunk in chunks_meta:
        data = restlib.findvalue(chunk['digest'])
        distrib_data = restlib.findvalue(chunk['digest'], local=False)
        if data != chunk['data']:
            print "DATA MISMATCH FOR DIGEST", chunk['digest']

        if distrib_data != chunk['data']:
            success = False
            print "DISTRIB DATA MISMATCH FOR DIGEST"
        
    if success:
        print len(chunks_meta), "items retrieved sucessfully"
finally:
    manager.kill_block_server(1234)
    manager.kill_block_server(1235)
