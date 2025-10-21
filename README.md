
# super-utils: Web Debugger Tool

super-utils is a portable, web-based debugging and network toolkit built with Go (Fiber) and HTMX. It allows you to run a wide range of system, network, and diagnostic tools from any browser, with a modern UI and easy deployment via Docker or Helm.

## Features

- **Run System & Network Tools:** Execute commands like `ping`, `tcpdump`, `nmap`, `curl`, `traceroute`, `iperf`, and many more directly from the web UI.
- **Curl & GraphQL:** Test HTTP endpoints and GraphQL queries interactively.
- **Network Info:** Lookup IP addresses, host details, and network information.
- **Single-File HTMX UI:** Lightweight, fast, and easy to use. The UI is embedded in the Go binary for portable execution.
- **Docker & Helm Support:** Deploy anywhere, from local dev to Kubernetes clusters.
- **Extensible:** Add more tools easily by updating the Dockerfile and UI dropdown.

## Supported Tools

super-utils includes a comprehensive set of CLI/network utilities (all available in the Docker image and UI):

### Network & Diagnostic Tools

apache2-utils, arping, bash, bind-tools, bird, bridge-utils, busybox, busybox-extras, calicoctl, conntrack-tools, ctop, curl, dhcping, drill, ethtool, file, fping, httpie, iftop, iperf, iperf3, iproute2, ipset, iptables, iptraf-ng, iputils, ipvsadm, liboping, mii-tool, mtr, net-snmp-tools, netcat-openbsd, netgen, nftables, ngrep, nmap, nmap-nping, scapy, socat, strace, tcpdump, tcptraceroute, termshark, tshark, websocat

### File & Text Processing

awk, cut, diff, find, grep, sed, vi, vim, wc, gzip, cpio, tar, jq

### System & Development

bash, git, libc6-compat, openssl, py-crypto, py2-virtualenv, python2, python3, util-linux

### Web & Transfer Tools

curl, wget, dig, nslookup, telnet, ssh, lftp, rsync, scp, ab (ApacheBench)

### Database Clients

mysql-client, postgresql-client, mongodb-tools, redis

## Quick Start

### Local Development

```bash
# Install dependencies
cd super-utils
go mod tidy
go run main.go
# Open http://localhost:8080 in your browser
```

### Docker

```bash
docker run -p 8080:8080 super-utils
# All tools are available inside the container

docker pull quay.io/pepodev/super-utils

# mysqldump
docker run --rm quay.io/pepodev/super-utils "mysqldump --help"

# gsutil // copy local file to google cloud storage
# working directory is /opt/ mount folder in this path to easy use as volume
docker run --rm -v ./service-account:./service-account quay.io/pepodev/super-utils "gcloud activate ./service-account/gcp && gsutils -m copy -r ./dir/ gs://some-bucket"
```

### Multi-Architecture Build (Local Testing)

To test multi-architecture builds locally, you need Docker Buildx:

```bash
# 1. Create a new builder instance (one-time setup)
docker buildx create --name multiarch --driver docker-container --use
docker buildx inspect --bootstrap

# 2. Build for multiple architectures
# Build and load for local testing (single platform at a time)
docker buildx build --platform linux/amd64 -t super-utils:amd64 --load .
docker buildx build --platform linux/arm64 -t super-utils:arm64 --load .

# 3. Test the built images
docker run --rm super-utils:amd64 /app/super-utils --version
docker run --rm super-utils:arm64 /app/super-utils --version

# 4. Build for multiple platforms and push to registry
docker buildx build --platform linux/amd64,linux/arm64 \
  -t quay.io/pepodev/super-utils:latest \
  --push .

# 5. Build without pushing (just verify it builds)
docker buildx build --platform linux/amd64,linux/arm64 \
  -t super-utils:multiarch .
```

**Note:**
- `--load` only works with single platform builds
- For multi-platform builds, you must use `--push` to a registry or omit both flags to just verify the build
- QEMU is automatically configured by Docker Desktop for cross-platform builds
- On Linux, you may need to install QEMU: `docker run --privileged --rm tonistiigi/binfmt --install all`

### Helm (Kubernetes)

```bash
cd super-utils/helm
helm install super-utils .
# Exposes service on port 8080
```

## Usage

1. **Open the Web UI:** Visit `http://localhost:8080` (or your server IP).
2. **Select a Tool:** Use the dropdown to choose a tool (e.g., `tcpdump`, `nmap`, `curl`).
3. **Enter Arguments:** Fill in the arguments field (see example hints for each tool).
4. **Run & View Output:** Click "Run Tool" to execute. Output appears below.
5. **Other Features:** Use the Curl/GraphQL and Network Info sections for specialized queries.

## API Endpoints

- `POST /api/shell` — Run selected tool (`tool`, `args`)
- `POST /api/curl` — Curl/GraphQL test (`url`, `graphql`)
- `POST /api/network` — Network info (`host`)

## Architecture

- **Go Fiber Backend:** Handles API requests, runs shell commands securely, and serves the embedded UI.
- **HTMX Frontend:** Single HTML file, embedded in the Go binary, provides a dynamic and responsive interface.
- **Dockerfile:** Installs all required tools for full functionality in containerized environments.
- **Helm Chart:** Simplifies deployment to Kubernetes clusters.

## File Structure

```text
super-utils/
├── main.go            # Fiber backend & API
├── web/index.html     # HTMX single-page UI (embedded)
├── Dockerfile         # Container build with all tools
├── helm/              # Helm chart for Kubernetes
│   ├── Chart.yaml
│   ├── values.yaml
│   └── templates/
│       ├── deployment.yaml
│       └── service.yaml
└── ...
```

## Security Notice

- This tool executes shell commands from the browser. Restrict access and run in trusted environments only.
- Consider using authentication, RBAC, or network restrictions for production deployments.
- Output is not sanitized—avoid running untrusted commands.

## License

MIT
