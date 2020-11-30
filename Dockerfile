FROM golang:1.14 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o linky-exporter cmd/linky-exporter/main.go


FROM alpine:3
COPY --from=builder /build/linky-exporter .
RUN addgroup -S exporter && adduser -S exporter -G exporter
USER exporter
ENTRYPOINT [ "./linky-exporter" ]
