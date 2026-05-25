# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /out/toolinfra ./cmd/toolinfra

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/toolinfra /toolinfra
USER nonroot:nonroot
EXPOSE 8080 8081
ENTRYPOINT ["/toolinfra"]
CMD ["serve", "--api", "--mcp-http"]
