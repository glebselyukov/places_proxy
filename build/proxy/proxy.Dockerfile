ARG APP_NAME=proxy

FROM golang:1.12-alpine3.9 AS builder
RUN apk add --update --no-cache make git bash
ARG APP_NAME
ARG SOURCE=/src/$APP_NAME
WORKDIR $SOURCE
COPY . .

ENV GO111MODULE=on \
    CGO_ENABLED=0
RUN go mod download
RUN go build -o /src/bin/$APP_NAME ./cmd/$APP_NAME/main.go

FROM alpine:3.9
ARG APP_NAME
COPY --from=builder /src/bin/$APP_NAME /usr/local/bin/$APP_NAME
ENV APP=$APP_NAME
ENTRYPOINT ["sh","-c", "$APP"]