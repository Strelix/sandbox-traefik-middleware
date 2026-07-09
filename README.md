# Strelix Traefik Middleware

A Traefik middleware plugin that logs request activity to Redis. It uses a `SET` command to Redis via the RESP protocol on every request, no redis client for simplicity.

This is used for Strelix Sandboxes and is how the controller queries activity.

## Configuration

### Plugin Config

| Name            | Type     | Description                                                           |
|:----------------|:---------|:----------------------------------------------------------------------|
| `redisAddr`     | `string` | The TCP address of the Redis server (e.g., `redis:6379`).             |
| `redisPassword` | `string` | (Optional) Password for Redis authentication.                         |
| `redisUser`     | `string` | (Optional) Username for Redis authentication (requires Redis 6+ ACL). |

## Installation

### As an Experimental Plugin (GitHub URL)

To use this plugin via Traefik's experimental plugins feature, add the following to your Traefik static configuration (`traefik.yml` or CLI):

```yaml
experimental:
  plugins:
    strelix-middleware:
      moduleName: "github.com/Strelix/sandbox-traefik-middleware"
      version: "v0.0.1"
```

### As a Local Plugin

If you want to run it locally without publishing to GitHub:

1. Clone this repository into a folder named `plugins-local/src/github.com/Strelix/sandbox-traefik-middleware`.
2. Configure Traefik to use the local plugin:

```yaml
experimental:
  localPlugins:
    strelix-middleware:
      moduleName: "github.com/Strelix/sandbox-traefik-middleware"
```

## Usage

Once the plugin is installed, you can use it in your dynamic configuration:

### File Provider (YAML)

```yaml
http:
  middlewares:
    my-redis-logger:
      plugin:
        strelix-middleware:
          redisAddr: "redis:6379"
          redisPassword: "optional-password"

  routers:
    my-router:
      rule: "Host(`example.com`)"
      service: my-service
      middlewares:
        - my-redis-logger

  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://localhost:8080"
```

### Docker Labels

```yaml
labels:
  - "traefik.http.middlewares.my-redis-logger.plugin.strelix-middleware.redisAddr=redis:6379"
  - "traefik.http.middlewares.my-redis-logger.plugin.strelix-middleware.redisPassword=optional-password"
  - "traefik.http.routers.my-router.middlewares=my-redis-logger"
```

## Development

A `dev/docker-compose.yml` is provided for local testing. It spins up Traefik and Redis and serves a basic nginx page.

```bash
cd dev
docker-compose up
```
