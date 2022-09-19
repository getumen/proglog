import requests
import base64

for i in range(3):
    requests.post(
        "http://localhost:8080",
        json={
            "record": {
                "value": base64.b64encode(f"Let's Go #{i}".encode()).decode(),
            }
        },
    )

for i in range(3):
    print(requests.get("http://localhost:8080", json={"offset": i}).content)
