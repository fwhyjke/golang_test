FROM golang:alpine AS builder

WORKDIR /usr/local/src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go test ./...
RUN CGO_ENABLED=0 go build -o ./bin/app ./cmd/app/


FROM alpine
COPY --from=builder /usr/local/src/bin/app /
CMD ["/app"]
