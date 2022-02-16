# Fuel-Tracker ![Build Status](https://github.com/roland-burke/fuel-tracker/actions/workflows/simple-workflow.yml/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/roland-burke/fuel-tracker/badge.svg?branch=master)](https://coveralls.io/github/roland-burke/fuel-tracker?branch=master)
This is the backend for the Fuel-Tracker App. It accepts basic HTTP requests and manages a Postgres Database. The setup is only tested on Linux.
## Prerequisite
For remote setup:
* running server
* docker installed and setup
* working ssh connection

For local setup:
* docker installed and setup
## Setup
First you should adjust the `config/conf.template.json` file e.g.:
```
{
    "apiKey": "<your_auth_token>",
    "port": 9006,
    "urlPrefix": "/fuel-tracker"
}
```

Then the environment varibles in `docker/.env.template` must be set e.g.:
```
PG_USERNAME=<my_postgres_username>
PG_PASSWORD=<my_pg_password>
PG_ADMIN_USER=<my_pgadmin_username>
PG_ADMIN_PASSWORD=<my_pgadmin_password>
```
**Important**: Rename all files with template in the name and replace template with prod or dev.

## How to run
To run the docker startup script, you should first create a docker-remote context:<br>
`docker context create <remote_name> ‐‐docker host=ssh://<user>@<remote_address>`

Edit the `deploy.sh` file and change the docker remote context name to <remote_name>

After that just run `deploy.sh prod` and you should be good to go.
For running in dev repeat the steps just for the files with dev extension instead of prod.

**Important**: Replace the config name which will be loaded in the Dockerfile
