# idp-app

A deliberately minimal Go HTTP service that serves as the sample workload for my
self-service **Internal Developer Platform (IDP)** on Kubernetes.

> **This app is intentionally trivial.** It is cargo, not the product.
> The point of the project is the platform around it: a `git push` here should
> end up as a scanned, secured, running deployment, with zero manual `kubectl`.

## How it fits in

The platform is split across two repositories:

| Repo | Responsibility |
|---|---|
| **idp-app** (this repo) | The application, its container image and the CI pipeline that builds, scans and publishes it |
| **idp-platform** *(coming soon)* | Infrastructure (Terraform) and Kubernetes manifests, reconciled into the cluster by ArgoCD |

This repo knows nothing about Kubernetes. Keeping "what the app is" separate from
"what should be running" mirrors real GitOps practice.

## Endpoints

| Method | Path | Response | Purpose |
|---|---|---|---|
| `GET` | `/health` | `{"status":"ok"}` | Liveness / readiness probes |
| `GET` | `/version` | `{"version":"<commit>"}` | Shows exactly which commit is running |

## Design choices

- **Go, standard library only:** zero dependencies, nothing to patch, no supply-chain risk.
- **Static binary:** (`CGO_ENABLED=0`) runs on an image with no OS libraries.
- **Multi-stage build on distroless:** the final image contains only the binary: no shell, no package manager.
- **Runs as non-root:** (`nonroot:nonroot`).
- **Version injected at build time:** via `-ldflags "-X main.version=..."` the commit SHA is baked into every build.
- **Configured via environment:** `PORT` (default `8888`).
- **`ReadHeaderTimeout` set:** protects against slow-header (Slowloris) attacks.

## Run it locally

Requirements: Go (see `go.mod` for the version), Docker, `make`.

```bash
# Run directly with Go (formats and vets first)
make run

# Build and run the container image
make docker-build
make docker-run            # serves on localhost:8888
make docker-run HOST_PORT=8888

# Try it
curl localhost:8888/health
curl localhost:8888/version
```

Configuration:

| Variable | Where | Default | Meaning |
|---|---|---|---|
| `PORT` | app (env) | `8888` | Port the app listens on |
| `VERSION` | make | `git describe --always --dirty` | Version stamped into the build |
| `HOST_PORT` | make | `8888` | Laptop port mapped to the container |

## CI pipeline *(in progress)*

On every push, GitHub Actions will:

1. Build the container image
2. Scan it with **Trivy** (the pipeline fails on critical vulnerabilities)
3. Push it to GHCR, tagged with the commit SHA
4. Update the image tag in `idp-platform`, where ArgoCD deploys it

## License

[MIT](LICENSE)