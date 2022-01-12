
## Compile the Go binary

FROM ghcr.io/invopop/golang:1.17.5-alpine AS build-go

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

ADD . /src
RUN go build -o gobl ./cmd/gobl

## Build Final Container

FROM alpine
RUN apk add --update --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=build-go /src/gobl /app/
COPY config/config.yaml /app/config/

VOLUME ["/app/config"]

# REST
EXPOSE 80
# gRPC
EXPOSE 8080

ENTRYPOINT [ "./gobl" ]
CMD [ "serve" ]
