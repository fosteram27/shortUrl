import requests

r = requests.get('http://localhost:8080/entries')

print(r.json())

