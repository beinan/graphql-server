# GraphQL Server Boilerplate in golang [WIP]


## Setup

We are using [Github](https://raw.githubusercontent.com/Masterminds/glide/) to manage the dependencies.

#### Install Glide 

```bash
curl https://glide.sh/get | shl
```




## Development

### Directory Structure

```
├── /server.go               # Http Service entry point
├── /schema/                 # GraphQL schema defination
├── /utils/                  # Utilities
│   ├── /password.go         # Password hashing via bcrypt
├── /vendor/                 # 3rd-party code managed by Glide
```

### Add new dependency

```bash
glide get github.com/{author}/{packange_name}
```

