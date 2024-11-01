FROM golang:latest AS build
WORKDIR /sensor_build
COPY . .
RUN CGO_ENABLED=0 go build -o sensor ./cmd/sensor

FROM alpine:latest
WORKDIR /sensor_run
COPY --from=build /sensor_build/sensor .

CMD ["./sensor"]
