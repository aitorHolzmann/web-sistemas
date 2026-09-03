FROM golang:1.27

RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    curl \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://github.com/sqlc-dev/sqlc/releases/download/v1.31.1/sqlc_1.31.1_linux_amd64.tar.gz \
    | tar -xz -C /usr/local/bin sqlc

WORKDIR /workspace