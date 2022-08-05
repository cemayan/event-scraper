# event-scraper

### Introduction

Event scraper is a event collector  which is getting data **Biletix** and **Passo**

In the future, it will be added  more provider

> **Biletix** is Turkey's leading ticketing company.
> **Passo** is also ticketing company in Turkey

It can be difficult to follow the events from every site. This app aims to close this gap.

---

### Requirements

For k8s:
- RabbitMQ [rabbitmq-operator](https://www.rabbitmq.com/kubernetes/operator/kubectl-plugin.html)
- PostgreSQL [kubegres](https://www.kubegres.io/doc/getting-started.html)
- Skaffold [skaffold](https://skaffold.dev/docs/install/)
- minikube etc [minikube](https://minikube.sigs.k8s.io/docs/start/)

### Usage

k8s:
```
skaffold dev
```
docker:
```
docker-compose up
```


### API

#### User Service

```http
POST :8089/api/v1/user
```
| Parameter  | Type | Description   |
|:-----------| :--- |:--------------|
| `username` | `string` | **Required**. |
| `password` | `string` | **Required**. |
| `email`    | `string` | **Optional**. |

Example Response:

```json
{
  "message": "User created {john.doe john.doe@test.com}"
}
```

#### Authorization

API requests require the API key.
To authenticate an API request, you should provide your API KEY in the **Authorization** header.

```http
POST :8109/api/v1/auth/getToken
```

| Parameter  | Type | Description   |
|:-----------| :--- |:--------------|
| `username` | `string` | **Required**. |
| `password` | `string` | **Required**. |

Example Response:

```json
{
    "message": "eyJhbG..."
}
```
#### Event API Service

Example Request:

```http
GET :8087/api/v1/event/provider/:provider
```

Example Response:

```json
[
  {
    "id": 26262,
    "Type": "MUSIC",
    "EventName": "*** ",
    "Place": "Jolly Joker Vadistanbul",
    "FirstDate": "2022-08-17 18:00:00 +0000 UTC",
    "SecondDate": "2022-08-17 18:00:00 +0000 UTC",
    "Provider": "BILETIX"
  },...
]
```

--- 
### Documentations

for  api: 
```
cd api 
godoc  -http=:6060 -notes=".*" -index  -goroot .
```

for  scraper:
```
cd scraper 
godoc  -http=:6060 -notes=".*" -index  -goroot .
```

for  user:
```
cd user 
godoc  -http=:6060 -notes=".*" -index  -goroot .
```

---

### Testing

In order to start test you should pass before the command

for api:
```
cd api
ENV="test" go test -v -cover -coverprofile=c.out  ./...
go tool cover -html=c.out   
```

for user:
```
cd user
ENV="test" go test -v -cover -coverprofile=c.out  ./...
go tool cover -html=c.out   
  
```
