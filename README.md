# Kikuri - きくり [WIP]

| Branch | Tests CI | Docker CD | Releases CD |
|--------|---------|-----------|-------------|
| `main` (stable)  | [![Unit Tests](https://github.com/zekurio/kikuri/actions/workflows/tests-ci.yml/badge.svg?branch=main)](https://github.com/zekurio/kikuri/actions/workflows/tests-ci.yml) | [![Docker CD](https://github.com/zekurio/kikuri/actions/workflows/docker-cd.yml/badge.svg?branch=main)](https://github.com/zekurio/kikuri/actions/workflows/docker-cd.yml) | [![Releases CD](https://github.com/zekurio/kikuri/actions/workflows/releases-cd.yml/badge.svg)](https://github.com/zekurio/kikuri/actions/workflows/releases-cd.yml) |
| `dev` (canary)   | [![Unit Tests](https://github.com/zekurio/kikuri/actions/workflows/tests-ci.yml/badge.svg?branch=dev)](https://github.com/zekurio/kikuri/actions/workflows/tests-ci.yml) | [![Docker CD](https://github.com/zekurio/kikuri/actions/workflows/docker-cd.yml/badge.svg?branch=dev)](https://github.com/zekurio/kikuri/actions/workflows/docker-cd.yml) |

## Inviting

You can choose between the stable and canary release, you can also use both at the same time. The stable release is the one that is recommended for most users, the canary release is the one that is recommended for developers and testers. The canary release is updated more frequently and may contain bugs and/or unfinished features.

[![INVITE STABLE](https://img.shields.io/badge/%20-INVITE%20STABLE-0288D1.svg?style=for-the-badge&logo=discord&color=5eac63)](https://discord.com/api/oauth2/authorize?client_id=1096747334987161652&permissions=285229072&scope=applications.commands%20bot)

[![INVITE CANARY](https://img.shields.io/badge/%20-INVITE%20CANARY-FFA726.svg?style=for-the-badge&logo=discord&color=3cfdd7)](https://discord.com/api/oauth2/authorize?client_id=1096748885734600829&permissions=285229072&scope=applications.commands%20bot)

## Features



## Self-hosting

If you want to self-host the bot, which is not recommended, you can do so by following the steps below.

### Using Docker

#### Use pre-built image

1. Install [Docker](https://docs.docker.com/get-docker/), on modern Linux distributions you can install it from your package manager, which should include `docker-compose` as well.
2. Create a `docker-compose.yml` file with the following content, you can orient yourself on the [example](https://github.com/zekurio/kikuri/blob/dev/config/example_docker-compose.yaml) file.
3. After you have created the `docker-compose.yml` file, you can start the bot by running `docker-compose up -d` in the same directory as the `docker-compose.yml` file.

#### Build image from source

1. Install [Docker](https://docs.docker.com/get-docker/), on modern Linux distributions you can install it from your package manager.
2. Clone the repository using `git clone https://github.com/zekurio/kikuri.git`, `cd` into the directory and build it with `docker build -t kikuri .`. This can take some time depending on your hardware.
3. Create a `docker-compose.yml` file with the following content, you can orient yourself on the [example](https://github.com/zekurio/kikuri/blob/dev/config/example_docker-compose.yaml)
4. After you have created the `docker-compose.yml` file, you can start the bot by running `docker-compose up -d` in the same directory as the `docker-compose.yml` file.

### Building from source



## Credits

* Icon from [Kou/tamagosan1001](https://twitter.com/tamagosan1001/status/1616869805520486400?s=20)
