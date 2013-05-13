import socket

SERVER = "localhost"
USE_EXTERNAL = False

# Taken from
# http://stackoverflow.com/questions/2311510/python-getting-a-machines-external-ip-address
def external_ip():
    sock = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
    sock.connect(('google.com', 80))
    ip = sock.getsockname()[0]
    sock.close()
    return ip

if USE_EXTERNAL:
    SERVER = str(external_ip())

print "Using server %s" % SERVER

