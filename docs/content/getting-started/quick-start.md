---
title: "Traefik Getting Started Quickly"
description: "Get started with Traefik Proxy and Docker."
---

# Quick Start

A Use Case Using Docker
{: .subtitle }

![quickstart-diagram](../assets/img/quickstart-diagram.png)

## Launch Traefik With the Docker Provider

Create a `docker-compose.yml` file where you will define a `reverse-proxy` service that uses the official Traefik image:

```yaml
services:
  reverse-proxy:
    # The official v3 Traefik docker image
    image: traefik:v3.4
    # Enables the web UI and tells Traefik to listen to docker
    command: --api.insecure=true --providers.docker
    ports:
      # The HTTP port
      - "80:80"
      # The Web UI (enabled by --api.insecure=true)
      - "8080:8080"
    volumes:
      # So that Traefik can listen to the Docker events
      - /var/run/docker.sock:/var/run/docker.sock
```

**That's it. Now you can launch Traefik!**

Start your `reverse-proxy` with the following command:

```shell
docker compose up -d reverse-proxy
```

You can open a browser and go to `http://localhost:8080/api/rawdata` to see Traefik's API rawdata (you'll go back there once you have launched a service in step 2).

## Traefik Detects New Services and Creates the Route for You

Now that you have a Traefik instance up and running, you will deploy new services.

Edit your `docker-compose.yml` file and add the following at the end of your file.

```yaml
services:

  ...

  whoami:
    # A container that exposes an API to show its IP address
    image: traefik/whoami
    labels:
      - "traefik.http.routers.whoami.rule=Host(`whoami.docker.localhost`)"
```

The above defines `whoami`: a web service that outputs information about the machine it is deployed on (its IP address, host, and others).

Start the `whoami` service with the following command:

```shell
docker compose up -d whoami
```

Browse `http://localhost:8080/api/rawdata` and see that Traefik has automatically detected the new container and updated its own configuration.

When Traefik detects new services, it creates the corresponding routes, so you can call them ... _let's see!_  (Here, you're using curl)

```shell
curl -H Host:whoami.docker.localhost http://127.0.0.1
```

_Shows the following output:_

```yaml
Hostname: a656c8ddca6c
IP: 172.27.0.3
#...
```

## More Instances? Traefik Load Balances Them

Run more instances of your `whoami` service with the following command:

```shell
docker compose up -d --scale whoami=2
```

Browse to `http://localhost:8080/api/rawdata` and see that Traefik has automatically detected the new instance of the container.

Finally, see that Traefik load-balances between the two instances of your service by running the following command twice:

```shell
curl -H Host:whoami.docker.localhost http://127.0.0.1
```

The output will show alternatively one of the following:

```yaml
Hostname: a656c8ddca6c
IP: 172.27.0.3
#...
```

```yaml
Hostname: s458f154e1f1
IP: 172.27.0.4
# ...
```

!!! question "Where to Go Next?"

    Now that you have a basic understanding of how Traefik can automatically create the routes to your services and load balance them, it is time to dive into [the user guides](../../user-guides/docker-compose/basic-example/ "Link to the user guides") and [the documentation](/ "Link to the docs landing page") and let Traefik work for you!

{!traefik-for-business-applications.md!}
