# Dockerfile for use as GitHub Action

# Builder
FROM cgr.dev/chainguard/go:latest as builder

WORKDIR /build
COPY main.go lib.go go.mod go.sum ./

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -v -o codeowners .

# Runner
FROM cgr.dev/chainguard/static

COPY --from=builder /build/codeowners /codeowners

ENTRYPOINT ["/codeowners"]