FROM golang:1.19-buster as builder

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -v -o migration ./migrations/
RUN go build -v -o server main.go

FROM debian:buster-slim

WORKDIR /app

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  ca-certificates && \
  rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/migrations/source /app/migrations/source
COPY --from=builder /app/server /app/server
COPY --from=builder /app/migration /app/migration
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/deploy/scripts/run.sh /app/run.sh

CMD ["/bin/bash","-c","./run.sh"]
