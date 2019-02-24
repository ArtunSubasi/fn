# ext-zeebe
[Fn](https://fnproject.io/) extension for [Zeebe.io](http://zeebe.io/)

**This is a prototype for a POC - it is not production ready!**

Features:
* Connects Fn functions to Zeebe to handle jobs. Each Fn function can be configured to handle a specific Zeebe job type.
* Starts Zeebe job workers which subscribe to the configured Zeebe job types and invoke the configured Fn functions
* Listens to the function deployment events to starts and stop the Zeebe job workers dynamically
* Provides a minimal REST endpoint to show the overview of registered Zeebe job workers

# Requirements
* [Zeebe](https://zeebe.io/) (tested using version 0.14.0)
* [Fn CLI](https://github.com/fnproject/cli) (testet using version 0.5.40)
* [Docker](https://www.docker.com/) (tested using Engine version 18.09.1)

# Fn Server modes
Fn server can be built in different modes: all-in-one, load balancer, API server and Fn runner. The following documentation describes the usage for the all-in-one mode.

# How to install the extension?
Fn extensions are currently installed by building a custom image of the Fn server including the extension. To build a custom image including the [Fn](https://fnproject.io/) extension for [Zeebe.io](http://zeebe.io/), make a file called `ext.yaml`.

```
extensions:
  - name: github.com/ArtunSubasi/ext-zeebe/zeebe
```

Then build a new custom docker image using the Fn CLI:

```
fn build-server -t imageuser/imagename
```

# Starting the Fn server with the extension

Start the server using docker.

```sh
docker run --rm -i --name fnserver \
    -e FN_LB_URL=http://localhost:8080 \
    -e FN_API_SERVER_URL=http://localhost:8080 \
    -e FN_ZEEBE_GATEWAY_URL=http://localhost:26500 \
    -v ./fn/data:/app/data  \
    -v /var/run/docker.sock:/var/run/docker.sock  \
    -p 8080:8080  \
    imageuser/imagename
```

## Environment variables
* FN_LB_URL: URL of the Fn Load Balancer. If using the all-in-one-mode, just point to the Fn server.
* FN_API_SERVER_URL: URL of the Fn API Server. If using the all-in-one-mode, just point to the Fn server.
* FN_ZEEBE_GATEWAY_URL: URL of the Zeebe Gateway (gRPC-Port)


## Docker-Volumes
* /app/data ist the database storing the deployed apps, functions, etc.
* /var/run/docker.sock points to the Unix-Socket of the Docker-Daemon so that the Fn server can manage internal docker containers. This is needed because the Fn functions are startet within their own docker containers.

# Configuring Fn functions to handle Zeebe jobs
Functions are configured to handle Zeebe jobs using the configuration parameter `zeebe_job_type` within the function configuration file `func.yaml`. An example of a function configuration:

```yaml
schema_version: 20180708
name: collectmoney
version: 0.0.4
runtime: go
entrypoint: ./func
format: http-stream
config:
  zeebe_job_type: payment-service
```
In the above example, the function `collectmoney` is configured to handle Zeebe jobs with the type `payment-service`. As soon as the function is deployed to the Fn, the extension launches Zeebe job workers and starts listening for available Zeebe jobs of the `payment-service`.

# Restrictions
Fn functions which are configured to handle Zeebe jobs must return a Json object as a response. The POC does not provide an automatic output mapping. Therefore, other return types, including Json arrays as a root, lead to an incident in the corresponding Zeebe workflow instance.
