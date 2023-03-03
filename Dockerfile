FROM golang:1.20-alpine as builder

WORKDIR /app
ARG VERSION

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg
RUN go build -ldflags="-X main.version=$VERSION" -o spbdayupdater ./cmd/spbdayupdater/main.go

FROM alpine:latest

RUN apk update && \
    apk add --no-cache tzdata

COPY --from=builder /app/spbdayupdater /app/spbdayupdater

CMD [ "/app/spbdayupdater" ]
