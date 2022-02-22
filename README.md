# integration hub gateway

integration-hub is the reverse proxy in the Integration Hub

## Getting Started

These instructions will get you a copy of the project up and running on your local machine.

### Prerequisites

* [Git](https://git-scm.com/) - Version control system

* [Docker Engine](https://store.docker.com/search?type=edition&offering=community) - Docker engine to manage Docker images and containers

* [Docker Compose](https://docs.docker.com/compose/) - Docker engine to run multiple containers from a docker-compose.yml file


### Installing

Make sure the OS meets all pre-requisites and then:

1. Install go

2. Add GOPATH to environment variable PATH
echo 'export PATH="/usr/local/opt/go@1.21/bin:$PATH"' >> ~/.zshrc (ZSH shell)

3. Clone repository
git@github.com:FJM-Slade/krakend-ce.git

https://github.com/FJM-Slade/krakend-ce.git

4. Execute script
./scripts/build_and_run.sh

```
$ cd ~ # this assumes project will be installed on home folder but any folder would be suitable
$ git clone https://github.com/FJM-Slade/krakend-ce.git
```

## Deploy in local environment

This assumes that all steps for the [Installing](#installing) section were done with success.

* Start application
You can use a new
```
$ docker-compose up -d
```


## Built with

* [KrakenD](https://www.krakend.io/) - API Gateway
