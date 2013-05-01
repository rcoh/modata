import json
import os

def create_new_keyfile():
    if not os.path.exists('key.kf'):
        with open('key.kf', 'w') as kf:
            kf.write('{}')
    return 'key.kf'

def save(name, content):
    with open(name, 'w') as kf:
        kf.write(json.dumps(content))
        
