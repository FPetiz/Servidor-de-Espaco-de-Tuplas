# ---------- Build ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o linda

# ---------- Runtime ----------
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/linda .

EXPOSE 54321

ENTRYPOINT ["./linda"]
