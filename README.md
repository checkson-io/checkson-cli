# checkson-cli

A command-line interface for [Checkson](https://checkson.io)

[![Build Status](https://github.com/checkson-io/checkson-cli/workflows/Lint%20%2F%20Test%20%2F%20IT/badge.svg?branch=main)](https://github.com/checkson-io/checkson-cli/actions)

[![asciicast](https://asciinema.org/a/sa694VnwgjjvxsvBIEecJWbEx.svg)](https://asciinema.org/a/sa694VnwgjjvxsvBIEecJWbEx)

## Installation

### Ubuntu / Debian

Download the newest .deb package for your architecture from the [releases page](https://github.com/checkson-io/checkson-cli/releases).

Then:

```
sudo apt install ./checkson-cli_1.0.16_linux_amd64.deb
```

### Snap

```
sudo snap install checkson
```

### Manual

Download the newest .tar.gz file for your architecture from the [releases page](https://github.com/checkson-io/checkson-cli/releases).

Then:

```bash
tar xvzf checkson-cli_1.0.16_linux_amd64.tar.gz -C /tmp
sudo mv /tmp/checkson /usr/local/bin
```

## Usage

### Login to Checkson

```bash
checkson login
```

You will be asked to login on https://app.checkson.io and to authorize the CLI.


### List checks

This shows the status of all checks:

```bash
checkson list
```

### Create check

This creates a new check that checks a website for SSL/TLS errors and sends
an email if a problem is found:

```bash
checkson-cli create new-check \
  --docker-image ghcr.io/checkson-io/checkson-testssl-check:main \
  --env URL=https://yourwebsite.com \
  --email me@example.com
```

### Show check details

This shows details of the given check:

```bash
checkson-cli show new-check
```

## Origin

This command line tool is partly based on [kafkactl](https://github.com/deviceinsight/kafkactl)

## Development

* Install [golangcli-lint](https://golangci-lint.run/usage/install/#local-installation)

## Releasing

In order to release a new version, do the following:

```
git tag -a v1.0.16 -m "v1.0.16"
git push origin v1.0.16
```
