import requests
import hashlib
SERVER = "http://localhost:1234/"
def store(chunk, digest):
    return requests.post(SERVER + "store", data={'key': digest, 'data': chunk}).json

def store_auto(chunk):
    return store(chunk, digest_for_chunk(chunk))


def findvalue(digest):
    resp = requests.get(SERVER + "find-value/" + digest).json
    if resp['status'] == 'OK':
        return resp['data']
    else:
        return None

def digest_for_chunk(chunk):
    return hashlib.sha1(chunk).hexdigest()
