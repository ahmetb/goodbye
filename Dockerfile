FROM golang:1.11-alpine as builder
RUN apk add --no-cache git
WORKDIR /src/goodbye
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN env CGO_ENABLED=0 go build -o /app ./cmd/goodbye

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /app /goodbye
ENTRYPOINT ["/goodbye"]
CMD ["-http-addr=ENV", "-followers-file=gs://ahmetb-goodbye/ids"]
