ARG TARGETOS
ARG TARGETARCH

FROM golang:1.22.2-alpine3.19 AS base

# See https://stackoverflow.com/a/55757473/4752298
ENV USER=appuser
ENV UID=12345

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/does_not_exist" \
    --no-create-home \
    --shell "/sbin/nologin" \
    --uid "$UID" \
    "${USER}"

WORKDIR /src
COPY go.* ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go mod download
COPY .. .

FROM base AS unit-test
RUN --mount=type=cache,target=/root/.cache/go-build \
    go test -v ./...

FROM base AS linter
ENV BINDIR=/usr/local/bin
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.62.2
RUN golangci-lint run


FROM base AS build-webserver
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-w -s" -o /out/webserver cmd/webserver/main.go

FROM base AS webserver
COPY --from=build-webserver /out/webserver /webserver
USER ${USER}:${USER}
ENTRYPOINT ["/webserver"]