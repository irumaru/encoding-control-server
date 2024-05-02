FROM golang:1.22 as app-builder

WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN go build -v -o app

FROM debian:stable
COPY --from=app-builder /app /usr/local/bin
CMD ["/usr/local/bin/app"]
