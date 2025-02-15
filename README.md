# timble

# Description

This service is take home technical test for a job application.

Timble is a simple Golang microservice for amobile dating app.

## Endpoints

See the endpoints served by timble in this [Postman Collection]()

## Technical Guidelines

### Requirements

Install these first, see the links for more details
1. [Git](https://git-scm.com/downloads)
2. [Golang 1.24+](https://go.dev/dl/)
3. [PostgreSQL](https://www.postgresql.org/download/)
4. [Redis](https://redis.io/docs/getting-started/installation/)

### Application Setup

#### Initial Setup

These steps only need to be done once.

1. Clone this repo. This repo already utilized go.mod, so you can clone it anywhere
```shell
git clone git@github.com:inassjunus/timble.git
```
2. Install dependencies
```shell
make prepare
```
3. Make sure redis and postgres already running.

```shell
# check for redis
redis-cli
# check for postgres
psql postgres
```

If they haven't, run them.
```shell
# starting redis in macOSX
brew services start redis
# starting redis in macOSX
brew services start postgresql
```
Check their respective guidelines to find the right command for your local machine.

4. Setup the database tables on postgres
```shell
# create admin user if you don't have one yet
CREATE ROLE timble WITH LOGIN PASSWORD <your password>;
ALTER ROLE timble CREATEDB;
```
Login with the admin user

```shell
psql postgres -U timble
```
Create database inside the psql CLI
```shell
CREATE DATABASE timble;
```
Execute the files from outside the CLI
```shell
psql -U timble -d timble -a -f db/migration/2025021313_create_users_table.sql
psql -U timble -d timble -a -f db/migration/2025021314_create_users_reactions_table.sql
```

5. Copy env.sample, then adjust the valus with the current environment details
```shell
cp env.sample .env
```
Please double check that the database & redis values are correct

#### Running the service

1. You can run with either executable file or with command

- Running with executable file

```shell
# Compile the package
make compile

# or, for MacOSX users:
make compile_osx

# make it executable
chmod +x deploy/_output/rest_timble

# run the executable file
./deploy/_output/rest_timble
```
Note: Make sure the `compile` command in the `Makefile` contains the correct os and arch for your local machine

- Running with `go run` command

```shell
make run-rest

```
2. To make sure the application running try running this command, it should return `ok`
```shell
curl localhost:9090//health
```

### Contribution
#### Linting
- Install [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) if you haven't
```shell
go install golang.org/x/tools/cmd/goimports@latest
goimports -l -w .
```
- Always ensure the code complies with the Golang syntax & import convention by running this command before committing

```shell
make pretty
```
#### Testing

Always run unit test before committing
```shell
make unit-test
```
Check for unit test coverage, we aim for 100% coverage
```shell
# display unit test coverage
make coverage
# display unit test coverage AND generate html file to check untested lines. The file will be names coverage.html, and you can open it on browser to see the untested lines
make coverhtml
```
