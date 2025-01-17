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

2. Open VS Code in the WSL:

   ```powershell
   code .
   ```

3. Clone the repo, and open VSCode in it:

   ```bash
   cd ~/

   git config --global user.name "Raki Rahman"
   git config --global user.email "mdrakiburrahman@gmail.com"

   git clone https://github.com/mdrakiburrahman/go-playground.git

   cd go-playground/
   code .
   ```

4. Fetch origin:

   ```bash
   git fetch origin
   ```

   Checkout any branch using VS Code UI.

5. Bootstrap your dev env

   ```bash
   GIT_ROOT=$(git rev-parse --show-toplevel)
   chmod +x ${GIT_ROOT}/contrib/bootstrap-dev-env.sh && ${GIT_ROOT}/contrib/bootstrap-dev-env.sh && source ~/.bashrc
   ```

Motes:

* If you run into docker problems, check `Docker Desktop: Settings > Resources > WSL Integration > Turn off/on Ubuntu-24.04`

## `hello-go` - a simple app

```bash
cd hello-go

go run main.go
go test ./...
```

The debugger settings should also work (first time debug takes a few seconds to boot):

![Debug Hello Go](./.imgs/debug-hello-go.png)

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