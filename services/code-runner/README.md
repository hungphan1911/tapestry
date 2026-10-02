# Code Runner

A self-hosted, synchronous competitive-programming executor. It provides a Codeforces-like GNU C++ / PyPy environment, not an identical Codeforces environment.

## Architecture

`POST /api/runs` is handled by the HTTP layer, validated by the run service, and dispatched through the `runner.Executor` interface. The MVP implementation invokes the Docker CLI, while the application layer has no Docker dependency. The language registry owns all filenames, images, compiler arguments, and interpreter commands; client input never becomes a shell command, executable name, Docker option, path, environment variable, or compiler flag.

## Languages

| ID | Image | Command |
| --- | --- | --- |
| `cpp` | `code-runner-cpp:2026.09` | `g++ main.cpp -std=c++23 -O2 -pipe -DONLINE_JUDGE -o main`, then `./main` |
| `python` | `code-runner-pypy:2026.09` | `pypy3 main.py` |

Build the execution images before starting the service:

```sh
docker build -t code-runner-cpp:2026.09 images/cpp
docker build -t code-runner-pypy:2026.09 images/python
docker run --rm code-runner-cpp:2026.09 g++ --version
docker run --rm code-runner-pypy:2026.09 pypy3 --version
```

The images are pinned to `gcc:14.2.0-bookworm` and `pypy:3.11-7.3.19-bookworm`; the commands above report the exact installed versions.

## Run

Docker Engine must be installed and running. On the host:

```sh
go run ./cmd/server
```

Or run the backend in Docker:

```sh
docker compose up --build
```

The compose deployment mounts the Docker socket only into the trusted backend so it can create execution containers. The backend never mounts that socket into execution containers.

Example:

```sh
curl -sS http://localhost:8080/api/runs \
  -H 'Content-Type: application/json' \
  -d '{"language":"cpp","source":"#include <iostream>\nint main(){int x;std::cin>>x;std::cout<<x*3<<"\\n";}","stdin":"5\n"}'
```

## Limits And Scheduling

The service waits for a slot while honoring the request context. The defaults are two concurrent runs, 512 KiB source, 1 MiB stdin, 256 MiB memory, one CPU, 64 PIDs, a 64 MiB workspace, 10-second compilation timeout, two-second runtime timeout, and 10 MiB each for stdout and stderr.

Responses use `success`, `compilation_error`, `runtime_error`, `time_limit_exceeded`, `output_limit_exceeded`, or `internal_error`. Docker's CLI does not provide a reliable portable distinction between OOM termination and other runtime failure, so this MVP intentionally does not report `memory_limit_exceeded`.

## Security Model

Each request receives a new container. It has no network, no host bind mounts, no Docker socket, runs as UID 10001, has a read-only root filesystem, uses a bounded tmpfs workspace, drops all Linux capabilities, sets `no-new-privileges`, and has memory/CPU/PID limits. Source is streamed directly into the workspace. Containers are removed after compilation errors, execution errors, timeouts, cancellation, and success.

Docker isolation is not a VM or a complete defense against kernel/container-runtime vulnerabilities. Keep Docker and the host patched, expose this service only over Tailscale as intended, and do not treat it as a multi-tenant public sandbox. The trusted backend's Docker-daemon access is a particularly sensitive boundary.

## Tests

```sh
go test ./...
go vet ./...
go test -tags=integration ./internal/executor
```

The integration suite requires Docker and the two execution images. It exercises C++ and PyPy success, compilation/syntax errors, runtime stderr, stdin/stdout, timeouts, output limits, outbound network denial, and inability to read a host-only temporary marker. PID exhaustion is environment-sensitive; manually submit the following PyPy source and verify it prints `limited` rather than `unbounded`:

```python
import os
children = []
try:
    for _ in range(1000):
        pid = os.fork()
        if pid == 0:
            os._exit(0)
        children.append(pid)
except OSError:
    print("limited")
else:
    print("unbounded")
for pid in children:
    os.waitpid(pid, 0)
```

The timeout and output-limit integration cases verify that programs cannot run or emit output indefinitely.

## Known Limitations

There is no authentication, persistence, queue, submission history, WebSocket API, problem storage, or exact Codeforces compatibility. C++ libraries, kernel behavior, compiler build flags, PyPy version, limits, and sandbox implementation differ from Codeforces. This service also relies on the local Docker CLI and daemon rather than a dedicated sandbox runtime.
