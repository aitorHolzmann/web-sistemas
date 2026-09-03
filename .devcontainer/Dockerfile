FROM golang:1.23

RUN apt-get update && apt-get install -y --no-install-recommends \
    RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1 \
    git \
    curl \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace
