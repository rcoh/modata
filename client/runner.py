import envoy
from multiprocessing import Process

class ServerManager(object):
    servers = {}
    def async_start_block_server(self, port):
        def inline():
            resp = envoy.run('go run ../main.go -block-server="localhost:%d"' % port)
            print resp.status_code
        p = Process(target=inline)
        p.start()
        self.servers[port] = p

    def kill_block_server(self, port):
        self.servers[port].terminate()

