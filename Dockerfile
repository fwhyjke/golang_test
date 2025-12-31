FROM golang:alpine AS builder

WORKDIR /usr/local/src

COPY go.mod ./

RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go test ./...
RUN go build -o ./bin/app ./cmd/app/

EXPOSE 8080

FROM alpine
COPY --from=builder /usr/local/src/bin/app /
CMD ["/app"]
