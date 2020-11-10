FROM golang:1.15.3 AS builder

COPY . /src/
WORKDIR /src/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o appexe .

FROM alpine:3.7
RUN apk --no-cache add ca-certificates

WORKDIR /bin/
COPY --from=builder /src/cmd/app .

EXPOSE 8080
ENTRYPOINT ["./app"]