# data-service

 The data-service exposes endpoints to create and read data.

## Requirements

* Docker >= 19.03
* Docker Compose >= 1.27.0
* Go >= 1.15

## Build

```shell
make build
```

## Test

```shell
make test
```

## Run

```shell
make run
```

## API

The following APIs are exposed by this service. 

1. Create `Data`

```shell
curl --request POST \
  --url http://localhost:8080/v1/data \
  --header 'Content-Type: application/json' \
  --data '{"title": "my_title"}'
```

2. Fetch `Data`

```shell
curl --request GET \
  --url http://localhost:8080/v1/data?title=my_title \
  --header 'Content-Type: application/json'
```

## Caveats

There is no unique constraint on the title, in case there are multiple records with the same title, the first created
is returned.
