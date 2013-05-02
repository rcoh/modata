import coding
import restlib
from runner import ServerManager

manager = ServerManager()

def verify_chunks(chunks_meta):
    for chunk in chunks_meta:
        distrib_data = restlib.findvalue(chunk['digest'], local=False)
        if distrib_data != chunk['data']:
            return False
    return True

try:
    n = 20
    ports = range(1000, 1000+n)
    for port in ports:
        manager.async_start_block_server(port)

    data_chunks = [str(x) for x in xrange(10)]
    chunks_meta = coding.make_chunk_dicts(data_chunks)
    
    for chunk in chunks_meta:
        # TODO: pick a random server
        restlib.store(chunk['data'], chunk['digest'], local=False)

    verify_chunks(chunks_meta)

    # Don't kill ports 0, who is us
    for i, port in enumerate(ports[1:]):
        chunks_all_good = verify_chunks(chunks_meta)
        manager.kill_block_server(port)
    
    
finally:
    for port in ports:
        manager.kill_block_server(port)
