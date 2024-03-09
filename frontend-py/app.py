from flask import Flask, request, render_template
import requests
import random

app = Flask(__name__)

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/about')
def about():
    return render_template('about.html')

@app.route('/url')
def url():
    result = ""
    status="Short_URL_Test"
    return render_template('url.html', status=status, result=result)

@app.route('/url', methods=['POST'])
def url_post():

    status = "processing..."
    text = request.form['urlLong']
    # post to server
    id = random.randint(1,100)
    idText = str(id)
    post = requests.post("http://localhost:8080/entries", json={"urlLong":text, "id":idText})
    urlShort = str(post.content.decode())

    return render_template('url.html', status=status, urlShort=urlShort)

@app.route('/shortUrl')
def shortUrl():
    r = requests.get('http://localhost:8080/entries')
    itemText = r.json()
    return render_template('shortUrl.html', result=itemText)

@app.route('/longUrl')
def longUrl():
    result = ""
    return render_template('longUrl.html', result=result)

@app.route('/longUrl', methods=['POST'])
def longUrl_post():

    # get URL to shorten from form
    text = request.form['urlLong']
    # post to server
    post = requests.post("http://localhost:8080/entries", json={"urlLong":text})
    urlShort = str(post.content.decode())

    return render_template('shortUrl.html', urlShort=urlShort)
