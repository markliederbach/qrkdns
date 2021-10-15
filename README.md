# QRKDNS <!-- omit in toc -->
This agent automatically discovers the current host's external IP address, and updates a given DNS A record in Cloudflare with the value.

- [Getting Started](#getting-started)
  - [Installation](#installation)
    - [Docker](#docker)

# Getting Started
## Installation
### Docker

Instead of installing directly, you can use the Docker images:

| Name | Description |
| ---- | ----------- |
| [ghcr.io/markliederbach/qrkdns](https://github.com/markliederbach/qrkdns/pkgs/container/qrkdns) | Basic image |

Example:

```console
docker run --rm --env-file .env ghcr.io/markliederbach/qrkdns
```
- The `.env` file should contain the following required environment variables
  - `NETWORK_ID` - The subdomain to update in Cloudflare (e.g., `myhost`)
  - `DOMAIN_NAME` - The base domain name (e.g., `foo.net`)
  - `CLOUDFLARE_ACCOUNT_ID` - Account ID from Cloudflare
  - `CLOUDFLARE_API_TOKEN` - Secret API token, with permission to read/update DNS records
