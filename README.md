# Fuel-Tracker
This is the backend for the Fuel-Tracker App.
## Setup
First you should adjust the `config/conf.template.json` file e.g.:
```
{
    "authToken": "<your_auth_token>",
    "port": 9006,
    "urlPrefix": "/fuel-tracker"
}
```
Important: Rename `config/conf.template.json` to `config/conf.json`

Then the environment varibles in `docker/.env.dev.template` must be set e.g.:
```
PG_USERNAME=<my_postgres_username>
PG_PASSWORD=<my_pg_password>
PG_ADMIN_USER=<my_pgadmin_username>
PG_ADMIN_PASSWORD=<my_pgadmin_password>
```
Important: Rename the file just like the conf.json file

## How to run
To run the docker startup script, you should first create a docker-remote context:<br>
`docker context create <remote-name> ‐‐docker host=ssh://<user>@<remoteaddress>`
