FROM golang:1.22.2-alpine3.18 as builder

WORKDIR /
COPY . 0g-da-op-plasma
RUN apk add --no-cache make
WORKDIR /0g-da-op-plasma
RUN make da-server

FROM alpine:3.18

COPY --from=builder /0g-da-op-plasma/bin/da-server /usr/local/bin/da-server

CMD ["da-server"]
