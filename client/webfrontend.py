from flask import Flask
from flask import request
import json
import requests
import coding
import time
from server_config import SERVER
import keyfilelib
from multiprocessing import Pool, Queue, Pipe, Process
from werkzeug import secure_filename


app = Flask(__name__)

fppage = """
<html>
<body>
<form name="upload a file" action="/upload" method="post">
<input type="filepicker" name="fileurl"/>
<input type="submit" value="Submit">
</form>
<script type="text/javascript" src="//api.filepicker.io/v1/filepicker.js"></script>
<script type="text/javascript">filepicker.setKey('AMIFT19ykQqibGJ2rxgdHz')</script>
</body>
</html>
"""

page = """
<!doctype html>
<html>
<body>
    <title>Upload new File</title>
    <h1>Upload new File</h1>
    <form action="/upload" method=post enctype=multipart/form-data>
      <p><input type=file name=file>
         <input type=submit value=Upload>
    </form>
</body>
<html>
"""

body_str = """
<html>
<body>
%s
</body>
</html>
"""

keyfile_name = keyfilelib.create_new_keyfile()
plength, clength = Pipe()
pname, cname = Pipe()

current_name = ""
current_size = 0 

@app.route("/fp")
def hello():
    return page

@app.route("/")
def index():
    global current_name
    global current_size
    global keyfile

    with open(keyfile_name, 'r') as keyfile_handle:
        keyfile = json.loads(keyfile_handle.read())

    if pname.poll():
        current_name = pname.recv()
        current_size = plength.recv()

    body = """
    <h3> MoData </h3>
    <div> Uploads
    <p><a href="/fp">Upload</a></p>
    <div> Waiting Uploads </div>
    <p> Current upload: %s </p>
    <p> %d pending </p>
    </div>
    <div> Downloads 
        <ul>
        %s
        </ul>
    </div>

    """ % (current_name + " - " + str(current_size) + " bytes",
           upload_jobs.qsize(),
           "".join(["<li><a href=/download/%s>Download %s</a></li>" % (key,key) for key in keyfile.keys()]))
    return body_str % body

@app.route("/upload", methods=['POST'])
def posted():
    if request.method == 'POST':
        ufile = request.files['file']
        if ufile:
            filename = secure_filename(ufile.filename)
            data = ufile.read()

    upload_jobs.put((filename, data))
    #metadata = coding.send_chunks_get_metadata(data)
    return body_str % """
    Queued up your file, waiting for server resources
    <p><a href="/">Watch Progress</a></p>
    """


def consume(input_queue, done_jobs, length_pipe, name_pipe, keyfile_name):
    while True:
        try:
            filename, data = input_queue.get()
            name_pipe.send(filename)
            length_pipe.send(len(data))

            metadata = coding.send_chunks_get_metadata(data)
            print "Got chunks"
            # Save it back
            with open(keyfile_name, 'r') as keyfile_handle:
                keyfile = json.loads(keyfile_handle.read())
            keyfile[filename] = metadata
            keyfilelib.save(keyfile_name, keyfile)
            print "Done saving to keyfile"
        except Queue.Empty:
            time.sleep(1)


if __name__ == "__main__":
    upload_jobs = Queue()
    done_jobs = Queue()

    app.debug = True
    consumer = Process(target=consume, args=(upload_jobs, done_jobs, clength, cname, keyfile_name))
    consumer.start()
    app.run(host=SERVER)
    p.join()
