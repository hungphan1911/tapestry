FROM golang:1.25.5-bookworm AS build

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /code-runner ./cmd/server

FROM docker:28.5-cli
COPY --from=build /code-runner /code-runner
EXPOSE 8080
ENTRYPOINT ["/code-runner"]
