# Proxy places

- [Proxy server](#project)
- [Sentry](#sentry)

---

## Project

#### Requirements

 * Golang 1.12+ _(optional 1.11.5+ for new go modules)_
 * Docker 1.10.0+
 * Docker Compose 1.6.0+ _(optional)_
 * Unix-like system (OS X | Linux) _(assuming you want to use the Makefile)_
 * [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) for `make format` cmd
 
 ...
 
 Some text
 
 ...

---

## Sentry

#### Requirements

 * Docker 1.10.0+
 * Docker Compose 1.6.0+ _(optional)_
 * Unix-like system (OS X | Linux) _(assuming you want to use the Makefile)_

#### Up and Running

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
