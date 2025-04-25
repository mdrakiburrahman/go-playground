# Go playground

## Style Guide

* [Uber's Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## Dev env setup

1. Get a fresh new WSL machine up:

   ```powershell
   # Delete old WSL
   wsl --unregister Ubuntu-24.04
   ```

   ```powershell
   # Create new WSL
   wsl --install -d Ubuntu-24.04
   ```

2. Clone the repo, and open VSCode in it:

   ```bash
   cd ~/

   git config --global user.name "Raki Rahman"
   git config --global user.email "mdrakiburrahman@gmail.com"

   git clone https://github.com/mdrakiburrahman/go-playground.git

   cd go-playground/
   code .
   ```

3. Reset your docker WSL integration since this is a new VM:

   > `Docker Desktop: Settings > Resources > WSL Integration > Turn off/on Ubuntu-24.04`

4. Bootstrap your dev env

   ```bash
   GIT_ROOT=$(git rev-parse --show-toplevel)
   chmod +x ${GIT_ROOT}/contrib/bootstrap-dev-env.sh && ${GIT_ROOT}/contrib/bootstrap-dev-env.sh && source ~/.bashrc
   ```

## `hello-go` - a simple app

```bash
cd hello-go

go run main.go
go test ./...
```

The debugger settings should also work (first time debug takes a few seconds to boot):

![Debug Hello Go](./.imgs/debug-hello-go.png)

## `sni-go` - a way to get an SNI cert

```bash
cd sni-go

az login --use-device-code

export VAULT_URL="https://a...d.vault.azure.net/"
export CERT_NAME="s...i"
export CLIENT_ID="e...7"
export TENANT_ID="7...7"
export SCOPE="https://database.windows.net/.default"

go run main.go
```

## `cert-auth-go` - use a local cert to auth

```bash
cd cert-auth-go

go run main.go \
   --cert-abs-path "${GIT_ROOT}/cert-auth-go/.secrets/myCert.cer" \
   --tenant-id "72f988bf-86f1-41af-91ab-2d7cd011db47" \
   --client-id "8b6e7cc1-d791-4844-bac4-cd50d649e63d" \
   --scope "https://management.core.windows.net/.default"
```

## OpenTelemetry

### Client/Server demo to Core Collector

Spin up the OTEL Collector (Core) binary in one terminal:

```bash
cd ${GIT_ROOT}
git clone https://github.com/open-telemetry/opentelemetry-collector.git
cd opentelemetry-collector
make install-tools
make otelcorecol
./bin/otelcorecol_* --config ./examples/local/otel-config.yaml
```

Spin up an HTTP Server that sends OTEL metrics to our Collector above:

```bash
cd ${GIT_ROOT}
git clone https://github.com/open-telemetry/opentelemetry-collector-contrib.git
cd opentelemetry-collector-contrib/examples/demo/server
go build -o server main.go; ./server
```

And finally, the HTTP Client that sends requests
```bash
cd ${GIT_ROOT}
cd opentelemetry-collector-contrib/examples/demo/client
go build -o client main.go; ./client
```

![Run Hello OTEL Core](./.imgs/run-otelcol-core-demo.png)

Debug the Client, Collector and Server together:

![Debug all three](./.imgs/debug-otelcol-core-all.png)

### Building a Custom Collector

Install `ocb` CLI:

```
curl --proto '=https' --tlsv1.2 -fL -o ocb \
https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.117.0/ocb_0.117.0_linux_amd64
chmod +x ocb
./ocb help
```

Build the collector:

```
./ocb --config custom-collector-builder-config.yaml
```

### Building a Receiver, Connector, Exporter

Spin up Jaeger UI:

```bash
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 14317:4317 \
  -p 14318:4318 \
  jaegertracing/all-in-one:1.41
```

Generate traces:

```
go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/telemetrygen@latest
telemetrygen traces --otlp-insecure --traces 1
```

View in Jaeger UI at `http://localhost:16686/`:

![Jaeger UI](./.imgs/jaeger-trace.png)

Initiate the go workspace:

```bash
cd ${GIT_ROOT}/opentelemetry-collector-raki

go work init
go work use otelcol-raki
go work use tailtracer
go work use exampleconnector
go work use emptyexporter
go work use marshaler
```

Run the Collector with the receiver wired up, either use VSCOde debugging, or via `go`:

```bash
go run ./otelcol-raki --config ${GIT_ROOT}/custom-collector-runtime-config-primary.yaml
```

* `custom-collector-runtime-config-primary.yaml`: Contains custom connectors
* `custom-collector-runtime-config-filelog.yaml`: Contains filelog receiver

> To add a specific package in go, e.g. we run `go get github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver@v0.115.0`

To demo the filelog receiver:

```bash
# Run in one terminal
go run ./otelcol-raki --config ${GIT_ROOT}/custom-collector-runtime-config-filelog.yaml

# Append in another terminal
export TEMP="/home/boor/go-playground/.temp"
export K8S_NAMESPACE="my-bar-namespace"
export K8S_POD="my-foo-pod"
export K8S_UID=$(uuidgen)
export K8S_CONTAINER="my-qux-container"
export K8S_CONTAINER_RESTART_NUM="0"
export FOLDER="${TEMP}/${K8S_NAMESPACE}_${K8S_POD}_${K8S_UID}/${K8S_CONTAINER}"
export FILE="${FOLDER}/${K8S_CONTAINER_RESTART_NUM}.log"

rm -rf ${TEMP}
mkdir -p ${FOLDER}
touch ${FILE}

echo "$(date '+%Y-%m-%d %H:%M:%S') ERROR This is a test error message" >> ${FILE}
echo "$(date '+%Y-%m-%d %H:%M:%S') DEBUG This is a test debug message" >> ${FILE}
echo "$(date '+%Y-%m-%d %H:%M:%S') INFO This is a test informational message" >> ${FILE}
```

The OTEL logs will show:

```
2025-03-01T13:34:09.156-0500    info    ResourceLog #0
Resource SchemaURL: 
Resource attributes:
     -> service.name: Str(my-bar-namespace/my-foo-pod/my-qux-container/0)
ScopeLogs #0
ScopeLogs SchemaURL: 
InstrumentationScope  
LogRecord #0
ObservedTimestamp: 2025-03-01 18:34:09.05617885 +0000 UTC
Timestamp: 1970-01-01 00:00:00 +0000 UTC
SeverityText: 
SeverityNumber: Unspecified(0)
Body: Str(2025-03-01 13:05:03 ERROR This is a test error message)
Attributes:
     -> container_name: Str(my-qux-container)
     -> restart_count: Str(0)
     -> namespace: Str(my-bar-namespace)
     -> log.file.path: Str(/home/boor/go-playground/.temp/my-bar-namespace_my-foo-pod_1e1f3d4b-c5a9-4092-a9c8-c7a4ea78c405/my-qux-container/0.log)
     -> log.file.name: Str(0.log)
     -> pod_name: Str(my-foo-pod)
     -> uid: Str(1e1f3d4b-c5a9-4092-a9c8-c7a4ea78c405)
Trace ID: 
Span ID: 
Flags: 0
        {"kind": "exporter", "data_type": "logs", "name": "debug"}
```

We don't do fancy RegEx parsing for timestamps and stuff. The `ObservedTimestamp` is good enough.