# Timble

# Description

This service is take home technical test for a job application.

Timble is a simple Golang microservice for a mobile dating app.

## Endpoints

See the endpoints served by Timble in this Postman collection
- [Timble.postman_collection.json](https://github.com/user-attachments/files/18813820/Timble.postman_collection.json)
- [Timble Local.postman_environment.json](https://github.com/user-attachments/files/18813824/Timble.Local.postman_environment.json)

## Technical Guidelines

### File structure

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
##### Linting initial setup
These steps only need to be done once.

- Install [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) if you haven't
```shell
go install golang.org/x/tools/cmd/goimports@latest
goimports -l -w .
```
- Install [staticcheck](https://github.com/dominikh/go-tools/tree/master/cmd/staticcheck)
```shell
go install honnef.co/go/tools/cmd/staticcheck@2025.1
```

##### Linting command
- Please always ensure the code complies with the Golang syntax & import convention by running this command before committing
```shell
make pretty
```
#### Testing

If you change any of the `interface`, rebuild the mocks files first with (Mockery)[https://vektra.github.io/mockery/latest/installation/] before running usit test
```shell
mockery --all --recursive --keeptree
```

Run basic unit test, this is the test that MUST be ran before each commit:
```shell
make unit-test
```

After significant changes, you need to run this too in order to check for race condition:
```shell
make race-test
```

If the unit tests are already passed, check for unit test coverage. We aim for 100% coverage.
```shell
# enable script fo be executed, this only need to be done once
chmod u+x ./tools/coverage.sh
# display unit test coverage
make coverage
# display unit test coverage AND generate html file to check untested lines. The file will be names coverage.html, and you can open it on browser to see the untested lines
make coverhtml
```
