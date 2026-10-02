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
- **Static binary** (`CGO_ENABLED=0`): runs on an image with no OS libraries.
- **Multi-stage build on distroless:** the final image contains only the binary, with no shell and no package manager.
- **Runs as non-root** (`nonroot:nonroot`).
- **Version injected at build time** via `-ldflags "-X main.version=..."`: the commit SHA is baked into every build.
- **Configured via environment:** `PORT` (default `8888`).
- **`ReadHeaderTimeout` set:** protects against slow-header (Slowloris) attacks.

## Run it locally

Requirements: Go (see `go.mod` for the version), Docker, `make`.

```bash
# Run directly with Go (formats and vets first)
make run                       # listens on localhost:8888
make run PORT=9090             # or any other port

# Build and run the container image
make docker-build
make docker-run                # serves on localhost:8888
make docker-run HOST_PORT=9090 # or any other laptop port

# Try it (container)
curl localhost:8888/health
curl localhost:8888/version
```

Configuration:

| Variable | Where | Default | Meaning |
|---|---|---|---|
| `PORT` | app (env) | `8888` | Port the app listens on (inside the container too) |
| `VERSION` | make | `git describe --always --dirty` | Version stamped into the build |
| `HOST_PORT` | make | `8888` | Laptop port mapped to the container's `8888` |

## CI pipeline

Every push and pull request runs [`ci.yml`](.github/workflows/ci.yml) on GitHub Actions:

1. **Check:** `gofmt` and `go vet` must pass
2. **Build:** container image, with the commit SHA injected as the version
3. **Scan:** Trivy fails the pipeline on fixable CRITICAL/HIGH vulnerabilities
4. **Smoke test:** the container must start, answer `/health`, and report the correct commit on `/version`
5. **Publish** *(pushes to `main` only)*: image pushed to GHCR, tagged with the full commit SHA

Pull requests go through steps 1–4 but are never published.

**Supply-chain hardening:** third-party actions are pinned by commit SHA (not tags),
the pipeline token is read-only by default and only the publishing job gets
`packages: write`, and no long-lived credentials are stored: GHCR login uses the
run's temporary `GITHUB_TOKEN`.

*Planned:* after publishing, update the image tag in `idp-platform` so ArgoCD deploys it.

## Container image

Published images: [`ghcr.io/vasileandreipopan/idp-app`](https://github.com/vasileandreipopan/idp-app/pkgs/container/idp-app),
one tag per commit on `main`, never `latest`.

```bash
docker run --rm -p 8888:8888 ghcr.io/vasileandreipopan/idp-app:<commit-sha>
curl localhost:8888/version
```

The app always listens on `8888` inside the container. Change only the left side
of `-p` (e.g. `-p 8888:8888`) to use a different port on your machine.

## License

[MIT](LICENSE)