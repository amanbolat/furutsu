FROM golang:1.15.3 AS builder

COPY . /src/
WORKDIR /src/cmd
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -mod=vendor -o appexe .

FROM alpine:3.7
RUN apk --no-cache add ca-certificates

WORKDIR /bin/
COPY --from=builder /src/cmd/appexe .
COPY migrates /migrates

EXPOSE 8080
ENTRYPOINT ["./appexe"]