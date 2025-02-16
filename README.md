# Timble

# Description

This service is take home technical test for a job application.

Timble is a simple Golang microservice for a mobile dating app.

## Endpoints

See the endpoints served by Timble in this [Postman](https://www.postman.com/) collection
- [Timble.postman_collection.json](https://github.com/user-attachments/files/18813820/Timble.postman_collection.json)
- [Timble Local.postman_environment.json](https://github.com/user-attachments/files/18813824/Timble.Local.postman_environment.json)

## Technical Guidelines

### File structure

#### Main directories
Our main codes are in these folders

<img width="176" alt="Screenshot 2025-02-16 at 14 32 09" src="https://github.com/user-attachments/assets/f327021a-3ca3-48a7-9612-352a54cbe354" />

- `cmd/`:  Contains the main applications for this project. Currently, there is only 1 application directory, which is `rest`. If there is any new one, such as `cronjob`, `grpc`, etc, they should be in separate directories under `cmd/`
- `internal/`: Contains shared codes that can be used by different applications inside the project.
     - `internal/config/`:  This is where we retrieve environment variables and perform application initial setup.
     - `internal/connection/`: Contains initial setup and direct calls for external connections (databases, other services, etc)
     - `internal/utils/`: Contains shared miscellaneous utility codes that are used throughout the project.
- `modules/`: Contains grouped logics for the service. Currently we only have `users` group. If there is any new one, such as `payment`, `chats`, etc, they should be in separate directories under `modules/`. The codes that directly interact with databases should be written here.
    - `modules/users/config/`: This directory contains initial setup for the group `users`
    - `modules/users/entity/`: This directory contains data structs for  `users`.
    - `modules/users/internal/`: We consider this a private folder, so it should not be imported outside `module/users/`
    - `modules/users/internal/handler`: Contains handler functions for each user-related endpoints
    - `modules/users/internal/repository`: Contains functions that call to other directories under `modules/` (there is none yet for now) and to `internal/connection/`
    - `modules/users/internal/usecase`: Contains main logic functions for each user-related endpoints

#### Additional directories
- `db/migrations`: Contains `.sql` files for database setup
- `mocks/`: Contains files for testing. This folder is autogenerated by mockery (see the [testing section](https://github.com/inassjunus/timble?tab=readme-ov-file#testing))
- `tools/`: Contains scripts to support development

### Requirements

Install these first, see the links for more details
1. [Git](https://git-scm.com/downloads)
2. [Golang 1.24+](https://go.dev/dl/)
3. [PostgreSQL](https://www.postgresql.org/download/)
4. [Redis](https://redis.io/docs/getting-started/installation/)

### Application Initial Setup

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

### Running the service

1. You can run with either executable file or with command

- Running with `go run` command (recommended)

```shell
make run-rest

```

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

2. To make sure the application running, try running this command; it should return response with a `ok` message
```shell
curl localhost:9090/health
```

3. Check the Prometheus metrics in browser; by default, it'll be in `localhost:8080`
<img width="1033" alt="Screenshot 2025-02-16 at 15 11 39" src="https://github.com/user-attachments/assets/5320b15d-04f9-4983-a5a0-792232956e8b" />


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

##### Unit Test

If you change any of the `interface`, rebuild the mocks files first with [Mockery](https://vektra.github.io/mockery/latest/installation/) before running unit test
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
# display unit test coverage AND generate html file to check untested lines.
make coverhtml
```

The generated file will be automatically named `coverage.html`, and you can open it on browser to see the untested lines

<img width="1173" alt="Screenshot 2025-02-16 at 15 17 12" src="https://github.com/user-attachments/assets/599a5c94-e74e-41cd-9e90-979b4d7b88f6" />

##### Postman Test
1. Import the Postman collection and environment from [here](https://github.com/inassjunus/timble?tab=readme-ov-file#endpoints) to you local Postman
2. Run the service based on [these steps](https://github.com/inassjunus/timble?tab=readme-ov-file#running-the-service)
3. Select a request and hit the `Send` button
4. Observe the endpoint response and test result
 <img width="1002" alt="Screenshot 2025-02-16 at 15 24 42" src="https://github.com/user-attachments/assets/93e569ae-a9f5-4e1b-8e56-bb723f426f1a" />
 <img width="994" alt="Screenshot 2025-02-16 at 15 24 36" src="https://github.com/user-attachments/assets/620b90a0-6c88-4205-bc11-465c19dc2f36" />


