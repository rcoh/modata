from flask import Flask
from flask import request
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
    print "file uploaded at: ", request.form['fileurl']
    return "got your file"


if __name__ == "__main__":
    app.debug = True
    app.run()
