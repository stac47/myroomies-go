# MyRoomies

This program aims at helping people living together in an houseshare.

It comes with a server part and a client part. The server exposes a simple HTTP
Restful API: hence the usage of client to access MyRoomies is not mandatory.

## Developement

### Dependencies

The program relies only on the Go Language: [Go language][go_lang_site]

Running the end-to-end tests will require you to install [jq][jq] on your
system:

```
brew install jq
```

### First Steps

Clone the repository using git:

```
$ git clone https://github.com/stac47/myroomies-go.git
```

Go into the fresh clones repository and update all the submodules:

```
$ git submodule update --init
```

## DevOps

### Installing the first time

The first time you start MyRoomies, you must provide an administrator password
with the environment variable `MYROOMIES_ROOT_PASSWORD`. By default, the
administrator login will be __root__: this can be changed with the environment
variable `MYROOMIES_ROOT_LOGIN`. These two variables are only useful on
MyRoomies first start.

```
MYROOMIES_ROOT_LOGIN=admin MYROOMIES_ROOT_PASSWORD=password docker-compose up -d
```

### Upgrading

Go into the git repository of MyRoomies, update it if needed with `git pull` or
select the tag you want to install `git fetch && git checkout v0.9.2`. Then,
deploy it:

```
docker-compose build myroomies-rest
docker-compose up --no-deps -d myroomies-rest
docker-compose logs myroomies-rest
```

[go_lang_site]: https://golang.org/
[jq]: https://stedolan.github.io/jq/
