# Stage 1: Build the Go binary
FROM golang:1.25.4-trixie AS builder

RUN adduser -uid 1111 builderuser

WORKDIR /app

COPY . .
COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o titler titler.go

# Stage 2: Minimal secure runtime
FROM scratch

WORKDIR /app

COPY --from=builder /app/titler .
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

USER 1111

EXPOSE 8080

ENTRYPOINT ["./titler"]