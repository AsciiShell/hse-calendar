FROM golang:1.12 AS builder
ADD . /go/src/github.com/asciishell/hse-calendar
WORKDIR /go/src/github.com/asciishell/hse-calendar
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/hse-calendar ./cmd/calendar-api/

FROM alpine:3.9
RUN apk add --update-cache ca-certificates
RUN apk add --no-cache tzdata
COPY --from=builder /go/src/github.com/asciishell/hse-calendar/bin/hse-calendar /usr/local/bin/hse-calendar
WORKDIR /usr/local/etc/hse-calendar
ENTRYPOINT ["/usr/local/bin/hse-calendar"]
