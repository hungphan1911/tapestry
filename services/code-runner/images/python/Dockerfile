FROM pypy:3.11-7.3.19-bookworm

RUN useradd --create-home --uid 10001 --shell /usr/sbin/nologin runner

USER runner
WORKDIR /workspace

# Inspect the installed runtime with: docker run --rm code-runner-pypy:2026.09 pypy3 --version
