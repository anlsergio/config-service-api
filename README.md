# Config Service API

## Development

### Prerequisites
- Go `v1.22+`
- Make: required by the helper scripts
- Docker
- Minikube
- Mockery: to generate and update mocks
- Swag CLI `v1.8.4`: to generate and update OpenAPI specs*
> [!WARNING]
> Swag must be installed in a specific version of `v1.8.4` because of [some issues](https://stackoverflow.com/questions/76582283/swag-init-generates-nothing-but-general-api-information) 
> recognizing annotations in dependency files.

### Getting started

Given the prerequisites are fulfilled, the only thing left to be able to run this project on your
local development environment is to prepare the Minikube cluster with some additional configurations.

#### Minikube Registry

For convenience, we are using the Minikube internal registry to push artifacts to when building from Docker
and to also pull images from, when deploying the K8s manifests.

```shell
minikube addons enable registry
```

#### Minikube Ingress
We have a dependency on the Kubernetes Ingress resource, so make sure to enable that too:

```shell
minikube addons enable ingress
```

#### Hosts file
Last but not least, update your local hosts file, so that your internal DNS translates the application host of
`config-service` into the Minikube cluster IP:

```shell
echo "$(minikube ip) config-service" | sudo tee -a /etc/hosts
```

### Helper Scripts

For commonly used tasks and commands, there are quite a few helper commands added to the `Makefile` 
of this project available for the `make` command, so make sure to check the file out to get to know the full list.

### Running the project

<details>

<summary>Running the application directly</summary>

Will spin up the application from your terminal
```shell
export SERVE_PORT=8080 && make run
```
> `SERVE_PORT` defines the port where the server will start listening for connections.

The application will be running at `localhost:8080`
```shell
curl http://localhost:8080/configs -v
```

</details>

<details>

<summary>Running the application from a Docker container</summary>

Will spin up the application container
```shell
make docker-up
```

The application will be running at `localhost:8080`
```shell
curl http://localhost:8080/configs -v
```

Update the docker container with your recent changes
```shell
make docker-update
```

</details>

<details>

<summary>Deploying the application to K8s</summary>

Deploy the application into your local Minikube instance
```shell
make deploy-k8s
```

</details>

### OpenAPI Documentation

Once the application is up and running, you should be able to access the Swagger endpoint, where the OpenAPI 
specifications for the routes implemented are parsed:

- If running from Minikube: http://config-service/swagger/index.html
- If running locally: http://localhost:8080/swagger/index.html

### Testing

Run the unit tests suite
```shell
make test
```

Run the unit tests suite while generating a coverage report
```shell
make test-cov
```

Render the test coverage report as HTML
```shell
make show-cov
```

[Back to top](#config-service-api)
