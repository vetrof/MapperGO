import requests

body = {
    "lat": "51.153705",
    "lng": "71.414391",
    "name": "bla bla bla"
}

response = requests.post('http://127.0.01:8080/create_place', json=body)

print(response.status_code)
print(response.content)
