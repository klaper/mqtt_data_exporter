FROM golang:1.13.1-alpine3.10 AS build

# Copy all project and build it
# This layer will be rebuilt when ever a file has changed in the project directory
WORKDIR /go/src/github.com/klaper_/mqtt_data_exporter
COPY . /go/src/github.com/klaper_/mqtt_data_exporter/
RUN go build -o /bin/mqtt_data_exporter

# This results in a single layer image
FROM alpine
COPY --from=build /bin/mqtt_data_exporter /bin/mqtt_data_exporter
EXPOSE 2112
ENTRYPOINT ["/bin/mqtt_data_exporter"]