FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o . ./...

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/quickrest .

EXPOSE 8090

ENTRYPOINT ["./quickrest"]

CMD ["--help"]