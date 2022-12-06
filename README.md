# knot-virtual-copergas

KNoT virtual is part of the KNoT project. It aims to provide an abstraction to allow certain protocols to interact with a cloud service using the KNoT AMQP Protocol, by virtualizing a KNoT Device. This service is responsible for virtualizing a Copergas API as a KNoT Device that retrieves data from Copergas API: https://gasonline.copergas.com.br/api/doc/index.html

## Contents

- [knot-virtual-copergas](#knot-virtual-copergas)
- [Basic installation and usage](#basic-installation-and-usage)
    - [System Installation and Usage](#system-installation-and-usage)
    - [Configuration](#configuration)
    - [Compiling and running](#compiling-and-running)
    - [Docker installation and usage](#docker-installation-and-usage)
      - [Docker requirements](#docker-requirements)
      - [Setting Environment Variables](#setting-environment-variables)
    - [Building and Running](#building-and-running)
      - [Production](#production)
      - [Development](#development)
  - [Verify service's health](#verify-services-health)
  - [Documentation](#documentation)
  - [License](#license)

# Basic installation and usage

### System Installation and Usage

- Go version 1.16+.
- Be sure the local packages binaries path is in the system's `PATH` environment variable:

```bash
# makes these following variables available to sub-processes from the current shell
export GOPATH=<path-to-golang> # in general, golang is installed at ${HOME}/go
export PATH=${PATH}:${GOPATH}/bin
```

### Compiling and running

```bash
make run
```

>You can use the `make watch` command to run the application on watching mode, allowing it to be restarted automatically when the code changes.

### Docker installation and usage

#### Docker requirements

To install the **project's Docker installation pre-requisites**, please follow the instructions in the link below:

- [Docker](https://docs.docker.com/get-docker/)

>_**Note**: if you're using a Linux system, please take a look at [Docker's post-installation steps for Linux](https://docs.docker.com/engine/install/linux-postinstall/)!_

#### Setting Environment Variables

Before proceeding, we must set the project's **environment variables** so that the service is properly configured. Environment variables will be set to the `.env` file and are ignored by `.gitignore` when committing changes. Under no circumstance one should push these files upstream.

To set the project's **environment variables**, change your current working directory to the project's root directory and create a `.env` file similar to the [.env.example](.env.example) one (make sure to update the values for all configurations):

```bash
# change current working directory
$ cd <path/to/knot-virtual-copergas>

# copies the content of .env.example to a new a file called .env
$ cp .env.example .env

# now, update the configuration values for each environment variable in .env
```

### Building and Running

Once you have all pre-requisites installed, change your current working directory to the project's root:

```bash
# change current working directory
$ cd <path/to/knot-virtual-copergas>
```

With Docker installed, to avoid problems caused by conflicts with the host's IP range, create a new Docker network. In the homologation environment, avoid networks with the following prefixes: 172.24, 172.25, and 172.26. This step needs to be done regardless of creating the image or downloading from the remote repository.

```bash
docker network create -d bridge --subnet=172.51.0.0/24 copergas
```

#### Production

In order to build the **Docker production image**, use the command below:

```bash
# build docker image from Dockerfile
$ docker build . --file Dockerfile --tag knot-virtual-copergas:latest
```

Finally, run the **production** container with the following command:

```bash
# start the container
$ sudo docker run -d --env-file .env --net=copergas -p 80:80 --restart unless-stopped --name knot-copergas -v $(pwd)/internal/config:/root/internal/config/ -it knot-virtual-copergas
```

If you want see log, use:

```bash
# create docker image
$ sudo docker logs -f <docker-image>
```

#### Development

In order to build the **Docker development image**, use the command below:

```bash
# build docker image from docker/Dockerfile-dev
$ docker build . --file docker/Dockerfile-dev --tag knot-virtual-copergas:dev
```

Finally, run the **development** container with the following command:

```bash
# start the container and clean up upon exit
$ docker run --rm --env-file .env --publish 8080:8080 --volume `pwd`:/usr/src/app --tty --interactive knot-virtual-copergas:dev
```

>_**Note**: the `--volume` flag binds and mounts `pwd` (your current working directory) to the container's `/usr/src/app` directory. This means that the changes you make outside the container will be reflected inside (and vice-versa). You may use your IDE to make code modifications, additions, deletions and so on, and these changes will be persisted both in and outside the container._


## Verify service's health

```bash
$ curl http://<hostname>:<port>/api/v1/healthcheck
```

## Documentation



Server documentation is auto-generated by the `swag` tool (<https://github.com/swaggo/swag>) from annotations placed in the code and can be viewed on the browser: `http://<address>:<port>/swagger/index.html`.

> If you want to generate the documentation just run the `make http-docs` command.

## License

Copyright © 2021-present, [CESAR](https://www.cesar.org.br). This project is [BSD-3-Clause](LICENSE) licensed.

