from flask import Flask
from flask import request
import requests
import coding
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

@app.route("/")
def hello():
    return page

@app.route("/upload", methods=['POST'])
def posted():
    fileurl = request.form['fileurl']
    data = requests.get(fileurl).text
    metadata = coding.send_chunks_get_metadata(data)
    return "got your file. Metadata: " + str(metadata)


if __name__ == "__main__":
    app.debug = True
    app.run()
