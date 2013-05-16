from flask import Flask
from flask import request, send_file, render_template, send_from_directory
import json
import coding
import time
from server_config import SERVER
import keyfilelib
from multiprocessing import Queue, Pipe, Process
from werkzeug import secure_filename
from cStringIO import StringIO


app = Flask(__name__)

keyfile_name = keyfilelib.create_new_keyfile()
plength, clength = Pipe()
pname, cname = Pipe()

current_name = ""
current_size = 0 

@app.route("/fp")
def upload():
    return render_template("upload.html")

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

    return render_template('index.html', current_upload=current_name, upload_size=current_size,
            num_pending=0, keys=sorted(keyfile.keys()))

@app.route("/upload", methods=['POST'])
def posted():
    if request.method == 'POST':
        ufile = request.files['file']
        if ufile:
            filename = secure_filename(ufile.filename)
            data = ufile.read()

    upload_jobs.put((filename, data))
    return render_template("queued.html")

@app.route("/download/<filename>", methods=['GET'])
def download(filename):
    with open(keyfile_name, 'r') as keyfile_handle:
        keyfile = json.loads(keyfile_handle.read())
    metadata = keyfile[filename]
    data = StringIO(coding.get_chunks(metadata))

    return send_file(data)

@app.route("/css/<filename>", methods=['GET'])
def css(filename):
    return send_from_directory(app.static_folder + "/css", filename)


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
    consumer.join()
