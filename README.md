# Proxy places

This service is a proxy of a third-party service that provides search sites. It provides fairly fast access
 to data through LRU caching algorithm, as well as updating data in the background after TTL expiration.


---

-
    - [Proxy server](#proxy-service)
        - [Requirements](#requirements-to-start-the-service)
        - [Environment variables](#environment-variables)
        - [Commands](#commands)
        - [Up and Running](#up-and-running)
    - [Sentry](#sentry)
        - [Requirements](#requirements-to-start-sentry-service)
        - [Up and Running](#up-and-running-sentry-service)

---

## Proxy service

This service provides information for cities, airports and countries based on one excellent avia service,
 from which we receive reliable information, normalize reponses and cache this data using the LRU algorithm.

### Requirements to start the service

 * Golang 1.12+ _(optional 1.11.5+ for new go modules)_
 * Docker 1.10.0+
 * Docker Compose 1.6.0+ _(optional)_
 * Unix-like system (OS X | Linux) _(assuming you want to use the Makefile)_
 * [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) for `make format` cmd

### Environment variables

|**Name**|**Description**|**Default value**|
|---|---|---|
|LOGGER_LEVEL|Level of logger (fatal, error, warning, info, debug)|info|
|LOGGER_STACKTRACE_LEVEL|Level of logger stacktrace (fatal, error, warning, info, debug)|error|
|LOGGER_SENTRY_LEVEL|Level of sentry logging (fatal, error, warning, info, debug)|error|
|LOGGER_SENTRY_DSN|Sentry DSN||
|LOGGER_SENTRY_STACKTRACE_ENABLED|Sentry stacktrace (true, false)|true|
|HTTP_SERVER_URL|Server address|:40001|
|HTTP_SERVER_NAME|Server name for sending in response headers|places_proxy|
|HTTP_SERVER_READ_TIMEOUT|Maximum duration in seconds for reading the full request (including body)|30|
|HTTP_SERVER_WRITE_TIMEOUT|Maximum duration in seconds for writing the full response (including body)|30|
|HTTP_CLIENT_NAME|Client name. Used in User-Agent request header|places_proxy|
|HTTP_CLIENT_READ_TIMEOUT|Maximum duration in seconds for full response reading (including body)|30|
|HTTP_CLIENT_WRITE_TIMEOUT|Maximum duration in seconds for full request writing (including body)|30|
|HTTP_CLIENT_UPDATE_MINUTES|Updating cached information by key after a certain number of minutes bypassing the LRU algorithm|10|
|HTTP_CLIENT_ENDPOINT|Endpoint of the third-party and excellent service to get places|`https://places.aviasales.ru/v2/places.json`|
|DB_URI|URI for connection to database|`redis://proxydefaultpass@127.0.0.1:50005/0`|

### Commands

|**Command**|**Description**|
|---|---|
|make prepare|Download packages and run generate commands|
|make format|Code formatting with `goimports` and `go fmt`|
|make project_build|Build project container and database container using `docker-compose`|
|make project_run|Running a pre-built project container and database container using `docker-compose`|
|make project_down|Stop and force remove project and database containers using `docker-compose`|

### Up and Running

With default environment
1. `make` or `go mod download` - Downloading packages.
2. `make project_build` - Build service and Redis containers.
3. `make project_run` - Run service and Redis, you are welcome.

Now you can make a GET request
As an example:

Request:

`GET: localhost:40001/api/places?locale=en&term=Moscow&types[]=city&types[]=airport`

Response:
```
[
    {
        "slug": "MOW",
        "subtitle": "Russia",
        "title": "Moscow"
    },
    {
        "slug": "DME",
        "subtitle": "Moscow",
        "title": "Moscow Domodedovo Airport"
    },
    {
        "slug": "SVO",
        "subtitle": "Moscow",
        "title": "Sheremetyevo International Airport"
    },
    {
        "slug": "VKO",
        "subtitle": "Moscow",
        "title": "Vnukovo Airport"
    },
    {
        "slug": "ZIA",
        "subtitle": "Moscow",
        "title": "Zhukovsky International Airport"
    },
    {
        "slug": "PUW",
        "subtitle": "United States",
        "title": "Pullman"
    }
]
```

The first request will pass as usual, and the second request for the same endpoint will already be cached.

On my computer, these are the results of this query:
- First request: 390ms
- Second request: 7ms

---

## Sentry

### Requirements to start Sentry service

 * Docker 1.10.0+
 * Docker Compose 1.6.0+ _(optional)_
 * Unix-like system (OS X | Linux) _(assuming you want to use the Makefile)_

### Up and Running Sentry service

Command #6 may take some time.

1. `cd ./build/sentry`
2. `docker volume create --name=sentry-data && docker volume create --name=sentry-postgres` - Make our local database and sentry volumes
    Docker volumes have to be created manually, as they are declared as external to be more durable.
3. `cp -n .env.example .env` - create env config file
4. `docker-compose build` - Build and tag the Docker services
5. `docker-compose run --rm web config generate-secret-key` - Generate a secret key.
    Add it to `.env` as `SENTRY_SECRET_KEY`.
6. `docker-compose run --rm web upgrade` - Build the database.
    Use the interactive prompts to create a user account.
7. `docker-compose up -d` - Lift all services (detached/background mode).
8. `cd ../..`
9. Access your instance at `localhost:9000`

If you walked through all the points correctly, then I congratulate you, you have a local Sentry running! :)
