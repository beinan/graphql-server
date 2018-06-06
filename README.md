# GraphQL Server Boilerplate in golang [WIP]


## Setup

We are using [Github](https://raw.githubusercontent.com/Masterminds/glide/) to manage the dependencies.

### Install Glide 

```bash
curl https://glide.sh/get | shl
```
### Install dependencies

```bash
glide install
```

### Install Docker Compose (Optional)


## Development

### Directory Structure

```
├── /server.go               # Http Service entry point
├── /schema/                 # GraphQL schema defination
├── /utils/                  # Utilities
│   ├── /password.go         # Password hashing via bcrypt
├── /vendor/                 # 3rd-party code managed by Glide
```

### Start a devel server via docker compose

With a single docker-compose command, we are able to start MongoDB, Redid and our graphql server. See Start devel server section.

#### Build the docker image of graphql server

Note: You only need to do a build when Dockerfile has changed
```bash
docker-compose build
```

#### Start devel server
```bash
docker-compose up
```
We are using [Fresh](https://github.com/pilu/fresh) for watching source files (*.go) changes, and restarting the graphql server inside a docker container.

You might need the setting below if you got "graphql_1  | inotify_init: too many open files"
```
fs.inotify.max_user_instances=32768
```
For arch linux users, you can put this into a config file e.g. /etc/sysctl.d/99-sysctl.conf


#### Stop devel server
Ctrl - C or 
```bash
docker-compose down
```


### Add new dependency

```bash
glide get github.com/{author}/{packange_name}
```

