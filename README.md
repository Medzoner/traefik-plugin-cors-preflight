# Traefik Plugin: CORS Preflight

[![Build Status](https://travis-ci.com/Medzoner/traefik-plugin-cors-preflight.svg?token=USx1h5scpzCMKrhJnEzv&branch=master)](https://travis-ci.com/github/Medzoner/traefik-plugin-cors-preflight)
[![Coverage Status](https://coveralls.io/repos/github/Medzoner/traefik-plugin-cors-preflight/badge.svg?branch=master&service=github)](https://coveralls.io/github/Medzoner/traefik-plugin-cors-preflight?branch=master)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go report](https://goreportcard.com/badge/github.com/Medzoner/traefik-plugin-cors-preflight?service=github)](https://goreportcard.com/report/github.com/Medzoner/traefik-plugin-cors-preflight?service=github)

## Short Description
Pass the browser cors preflight with response status 204 for Method OPTIONS

## Configuration


Requirements: Traefik >= v2.5.5

### Static

```bash
--pilot.token=xxx
--experimental.plugins.corspreflight.modulename=github.com/Medzoner/traefik-plugin-cors-preflight
--experimental.plugins.corspreflight.version=v1.0.7
```

```yaml
pilot:
  token: xxx

experimental:
  plugins:
    corspreflight:
      modulename: github.com/Medzoner/traefik-plugin-cors-preflight
      version: v1.0.7
```

```toml
[pilot]
    token = "xxx"

[experimental.plugins.corspreflight]
    modulename = "github.com/Medzoner/traefik-plugin-cors-preflight"
    version = "v1.0.7"
```

```yml
testData:
  testData:
    code: 204
    method: 'OPTIONS'
```

### Dynamic

To configure the `CorsPreflight` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/).

#### File 

```yaml
http:
  middlewares:
    corspreflight-middleware:
      plugin:
        corspreflight:
          code: 200
          method: OPTIONS

  routers:
    my-router:
      rule: Host(`localhost`)
      middlewares:
        - corspreflight-middleware
      service: my-service

  services:
    my-service:
      loadBalancer:
        servers:
          - url: 'http://127.0.0.1'
```

```toml
[http.middlewares]
  [http.middlewares.corspreflight-middleware.plugin.corspreflight]
    code = 200
    method = "OPTIONS"

[http.routers]
  [http.routers.my-router]
    rule = "Host(`localhost`)"
    middlewares = ["corspreflight-middleware"]
    service = "my-service"

[http.services]
  [http.services.my-service]
    [http.services.my-service.loadBalancer]
      [[http.services.my-service.loadBalancer.servers]]
        url = "http://127.0.0.1"
```

#### Kubernetes

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: corspreflight-middleware
spec:
  plugin:
    corspreflight:
      code: 200
      method: OPTIONS

---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: whoami
spec:
  entryPoints:
    - web
  routes:
    - kind: Rule
      match: Host(`whoami.localhost`)
      middlewares:
        - name: corspreflight-middleware
      services:
        - kind: Service
          name: whoami-svc
          port: 80
```

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: corspreflight-middleware
spec:
  plugin:
    corspreflight:
      code: 200
      method: OPTIONS

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myingress
  annotations:
    traefik.ingress.kubernetes.io/router.middlewares: default-corspreflight-middleware@kubernetescrd

spec:
  rules:
    - host: whoami.localhost
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name:  whoami
                port:
                  number: 80
```

#### Docker

```yaml
version: '3.7'

services:
  whoami:
    image: traefik/whoami:v1.7.1
    labels:
      traefik.enable: 'true'

      traefik.http.routers.app.rule: Host(`whoami.localhost`)
      traefik.http.routers.app.entrypoints: websecure
      traefik.http.routers.app.middlewares: corspreflight-middleware
      
      traefik.http.middlewares.corspreflight-middleware.plugin.corspreflight.code: 200
      traefik.http.middlewares.corspreflight-middleware.plugin.corspreflight.method: 'OPTIONS'
```

## Developed & Maintained by
[Mehdi Youb](https://github.com/Medzoner)

## License
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)