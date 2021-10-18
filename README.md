# QRKDNS <!-- omit in toc -->
This agent automatically discovers the current host's external IP address, and updates a given DNS A record in Cloudflare with the value.

- [Getting Started](#getting-started)
  - [Installation](#installation)
    - [Docker](#docker)
      - [Examples](#examples)

# Getting Started
## Installation
### Docker

Instead of installing directly, you can use the Docker images:

| Name | Description |
| ---- | ----------- |
| [ghcr.io/markliederbach/qrkdns](https://github.com/markliederbach/qrkdns/pkgs/container/qrkdns) | Basic image |

**A note about `.env` files and docker.** The `--env-file` option only takes env files that do not have quotations. If quotes are present, they won't parse the options correctly. It is recommended that you create a separate file called `.env.docker`, which doesn't use any quotations.


#### Examples

**Sync Once**
```console
docker run --env-file .env.docker --rm -it  ghcr.io/markliederbach/qrkdns:latest sync
```
- The `.env.docker` file should contain the following required environment variables
  - `NETWORK_ID` - The subdomain to update in Cloudflare (e.g., `myhost`)
  - `DOMAIN_NAME` - The base domain name (e.g., `foo.net`)
  - `CLOUDFLARE_ACCOUNT_ID` - Account ID from Cloudflare
  - `CLOUDFLARE_API_TOKEN` - Secret API token, with permission to read/update DNS records

**Sync Cron**
```console
docker run --env-file .env.docker --rm -it  ghcr.io/markliederbach/qrkdns:latest sync cron
```
- The `.env.docker` file should contain the following required environment variables
  - `NETWORK_ID` - The subdomain to update in Cloudflare (e.g., `myhost`)
  - `DOMAIN_NAME` - The base domain name (e.g., `foo.net`)
  - `CLOUDFLARE_ACCOUNT_ID` - Account ID from Cloudflare
  - `CLOUDFLARE_API_TOKEN` - Secret API token, with permission to read/update DNS records
  - `SCHEDULE` - Cron pattern describing how often the sync job should be run
