FROM golang:1.16-alpine AS build

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir -p /src

WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download

COPY . /src
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o streamer

FROM alpine:3.11

RUN apk add --update tzdata ca-certificates bash && \
    mkdir -p /app && \
    chgrp -R 0 /app && \
    chmod -R g=u /app

WORKDIR /app

COPY --from=build /src/streamer /app

CMD ["./streamer"]
