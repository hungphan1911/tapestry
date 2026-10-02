FROM gcc:14.2.0-bookworm

RUN useradd --create-home --uid 10001 --shell /usr/sbin/nologin runner

USER runner
WORKDIR /workspace

# Inspect the installed compiler with: docker run --rm code-runner-cpp:2026.09 g++ --version
