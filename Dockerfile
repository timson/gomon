## Build
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY plugins ./plugins

RUN go build -o /gomon

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /gomon /gomon

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/gomon serve"]