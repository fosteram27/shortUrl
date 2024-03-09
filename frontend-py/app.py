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

    status = "processing..."
    text = request.form['text']
    # post to server
    id = random.randint(1,100)
    idText = str(id)
    post = requests.post("http://localhost:8080/entries", json={"urlLong":text, "id":idText})

    # fieldId = 'entry' + idText
    getText = 'http://localhost:8080/entries/' + 'entry' + idText
    # getText = 'http://localhost:8080/entries/' + fieldId
     
    r = requests.get(getText)
    
    itemText = r.json()
    return render_template('shortUrl.html', status=status, result=itemText)
