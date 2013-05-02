import envoy
import requests
from multiprocessing import Process
import time

def check_live(port):
    live = False
    while not live:
        try:
            resp = requests.get("http://localhost:%d/ping" % port)
            print resp
            live = True
        except:
            time.sleep(.1)

def check_dead(port):
    live = True
    while live:
        try:
            requests.get("http://localhost:%d/ping" % port)
            time.sleep(.1)
        except:
            live = False
            print "Dead"


class ServerManager(object):
    servers = {}
    def async_start_block_server(self, port):
        def inline():
            resp = envoy.run('go run ../main.go -block-server="localhost:%d"' % port)
            print resp.std_out
        p = Process(target=inline)
        p.start()
        check_live(port)
        self.servers[port] = p

    def kill_block_server(self, port):
        self.servers[port].terminate()
        check_dead(port)


