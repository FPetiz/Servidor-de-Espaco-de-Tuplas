# ---------- Build ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
# Se você tiver um go.sum, adicione aqui também
RUN go mod download

COPY . .

# Usamos uma variável de argumento para definir a pasta (mutex ou channels)
ARG VERSION=mutex

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o linda ./${VERSION}

# ---------- Runtime ----------
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/linda .

EXPOSE 54321

ENTRYPOINT ["./linda"]