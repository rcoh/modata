import envoy
import requests
from multiprocessing import Process
import time
import server_config

SERVER = server_config.SERVER

def check_live(port):
    live = False
    print "Waiting "
    while not live:
        try:
            resp = requests.get("http://%s:%d/ping" % (SERVER, port))
            print resp
            live = True
        except:
            print ".",
            time.sleep(.1)
    print "Done"

def check_dead(port):
    live = True
    while live:
        try:
            requests.get("http://%s:%d/ping" %(SERVER, port))
            time.sleep(.1)
        except:
            live = False
            print "Dead"


class ServerManager(object):
    servers = {}
    def async_start_block_server(self, port, replication_port=None):
        def inline():
            if self.servers:
                bootstrap = '-bootstrap="%s:%d"' % (SERVER, self.servers.keys()[0])
            else:
                bootstrap = ""

            if replication_port is None:
                resp = envoy.run('go run ../main.go -block-server="%s:%d" %s' % (server, port, bootstrap))
            else:
                resp = envoy.run('go run ../main.go \
                        -block-server="%s:%d" %s -replication=true \
                        -replication-server="%s:%d"' % (SERVER, port, bootstrap, SERVER, replication_port))

            print resp.std_out
        p = Process(target=inline)
        p.start()
        check_live(port)
        self.servers[port] = p

    def kill_block_server(self, port):
        self.servers[port].terminate()
        check_dead(port)

    def kill_all(self):
        for port in self.servers:
            self.servers[port].terminate()
            check_dead(port)

