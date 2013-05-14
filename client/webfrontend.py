from flask import Flask
from flask import request
import json
import requests
import coding
from server_config import SERVER
import keyfilelib
app = Flask(__name__)

page = """
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

body_str = """
<html>
<body>
%s
</body>
</html>
"""

keyfile_name = keyfilelib.create_new_keyfile()
with open(keyfile_name, 'r') as keyfile_handle:
    keyfile = json.loads(keyfile_handle.read())


@app.route("/fp")
def hello():
    return page

@app.route("/")
def index():
    body = """
    <div> MoData
    <a href="/fp">Upload</a>
    <ul>
    %s
    </ul>
    </div>

    """ % "".join(["<li><a href=/download/%s>Download %s</a></li>" % (key,key) for key in keyfile.keys()])
    return body_str % body

@app.route("/upload", methods=['POST'])
def posted():
    fileurl = request.form['fileurl']
    data = requests.get(fileurl).text
    metadata = coding.send_chunks_get_metadata(data)
    return "got your file. Metadata: " + str(metadata)


if __name__ == "__main__":
    app.debug = True
    app.run(host=SERVER)
