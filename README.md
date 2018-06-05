This is an experimental app written in Go on the server and client.

# Prerequisites

 - [Go compiler](https://golang.org)
 - [Docker](https://docker.com)
 - [gcloud](https://cloud.google.com/sdk/gcloud/)

# Server

## Build

From the repo root, build the server container and give it a name:

```bash
$ SERVER_NAME=server
$ docker build -f server/Dockerfile -t $SERVER_NAME .
```

## Local

Run the server container and map ports.

```bash
$ docker run --rm -p 8080:8080 $SERVER_NAME
```

Open [http://localhost:8080](http://localhost:8080) to view the app..

Use CTRL-C to terminate the server.

## Cloud (GKE)

Configure gcloud for your project and install the Kubernetes command line tool:

```bash
$ gcloud config set project [PROJECT_ID]
$ gcloud config set compute/zone us-central1-b
$ gcloud components install kubectl
$ gcloud auth configure-docker
```

Create a cluster:

```bash
$ gcloud container clusters create ${SERVER_NAME}-cluster \
  --num-nodes=3 \
  --cluster-version=1.10.2-gke.3 \
  --machine-type f1-micro \
  --addons HorizontalPodAutoscaling,HttpLoadBalancing \
  --enable-autoupgrade \
  --enable-autorepair \
  --scopes "https://www.googleapis.com/auth/devstorage.read_only", \
           "https://www.googleapis.com/auth/logging.write", \
           "https://www.googleapis.com/auth/monitoring", \
           "https://www.googleapis.com/auth/servicecontrol", \
           "https://www.googleapis.com/auth/service.management.readonly", \
           "https://www.googleapis.com/auth/trace.append"
```

> **_Optional:_** You can also create a cluster [via the Console](https://console.cloud.google.com/kubernetes/list) or use a previously created cluster.
>
> To obtain credentials for an existing cluster or a cluster not created using `gcloud`:
>
> ```bash
> $ gcloud container clusters get-credentials ${SERVER_NAME}-cluster
> ```

Tag the image (built above), push it and run it:

```bash
$ export PROJECT_ID="$(gcloud config get-value project -q)"
$ docker tag ${SERVER_NAME} gcr.io/${PROJECT_ID}/${SERVER_NAME}
$ docker push gcr.io/${PROJECT_ID}/${SERVER_NAME}:latest
$ kubectl run ${SERVER_NAME} --image=gcr.io/${PROJECT_ID}/${SERVER_NAME} --port 8080
```

> **_Optional_**: create a static IP (for domain mapping):
> ```bash
> $ gcloud compute addresses create ${SERVER_NAME}-ip --region > us-central1
> $ CLUSTER_IP=$(gcloud compute addresses describe $SERVER_NAME-ip --region us-central1 --format 'value(address)')
> ```

Expose the app to the internet:

```bash
$ kubectl expose deployment ${SERVER_NAME} --type=LoadBalancer --port 80 --target-port 8080
$ CLUSTER_IP=$(kubectl get service --field-selector=metadata.name=${SERVER_NAME} --output=custom-columns=ip:status.loadBalancer.ingress[0].ip --no-headers)
```

> **_Optional_**: if using a static IP:
> ```bash
> $ kubectl expose deployment ${SERVER_NAME} --type=LoadBalancer --port 80 --target-port 8080 --load-balancer-ip=$CLUSTER_IP
> ```

Navigate to the app:

```bash
$ xdg-open "http://$CLUSTER_IP"
```

To update the app, build and push a new image:

```bash
$ docker build -f server/Dockerfile -t $SERVER_NAME .
$ docker tag ${SERVER_NAME} gcr.io/${PROJECT_ID}/${SERVER_NAME}
$ docker push gcr.io/${PROJECT_ID}/${SERVER_NAME}
$ kubectl set image deployment/${SERVER_NAME} ${SERVER_NAME}=gcr.io/${PROJECT_ID}/${SERVER_NAME}
```

# Android

From the repo root, build the Android container (which will also build the app):

```bash
$ docker build -f android/Dockerfile -t android .
```

Run the container to copy the outputs. The following places the outputs where Android Studio would have:

```bash
$ OUTPUT_ABS_PATH=${PWD}/android/app/build/outputs
$ mkdir -p $OUTPUT_ABS_PATH
$ docker run --rm --user $UID:$(id -g) -v ${OUTPUT_ABS_PATH}:/outputs android
```

# iOS

TBD

# IDE

For IDE auto-completion, install the go-wasm compiler to your system. You will need a regular [Go compiler](https://golang.org/) installed on your system and the go-wasm compiler will live alongside it.

```bash
$ git clone --branch wasm-wip https://github.com/neelance/go.git $HOME/go-wasm
$ cd $HOME/go-wasm/src && ./all.bash
```

To use the WASM-capable compiler for the duration of a terminal session:

```bash
$ GOROOT="$HOME/go-wasm"
$ alias go="$HOME/go-wasm/bin/go"
```

Configure VS Code by adding the following to your workspace settings (you'll need to expand out $HOME yourself):

```json
"go.goroot": "$HOME/go-wasm"
```

# Resources

 - https://github.com/neelance/go
 - https://blog.lazyhacker.com/2018/05/webassembly-wasm-with-go.html
