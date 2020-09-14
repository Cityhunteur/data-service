# data-service

 The data-service exposes endpoints to create and read data.

## Requirements

* Docker >= 19.03
* Docker Compose >= 1.27.0
* Go >= 1.15

## High Level Design


## Build

```
make build
```


## Test

```
make test
```

## Run

```
make run
```

## API

The following APIs are exposed by this service. 

1. Create `Data`

```
curl --request POST \
  --url http://localhost:8080/v1/data \
  --header 'Content-Type: application/json' \
  --data '{"title": "my_title"}'
```

2. Fetch `Data`

```
curl --request GET \
  --url http://localhost:8080/v1/data?title=my_title \
  --header 'Content-Type: application/json'
```


## Caveats

There is no unique constraint on the title, in case there are multiple records with the same title, the first created
is returned.