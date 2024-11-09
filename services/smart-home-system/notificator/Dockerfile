FROM golang:latest AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd/notificator

FROM alpine:latest
WORKDIR /run
COPY --from=build /build/app .

CMD ["./app"]
