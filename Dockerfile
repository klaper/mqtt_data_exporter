FROM golang:1.12.6-alpine3.9 AS build

# Install tools required to build the project
# We will need to run `docker build --no-cache .` to update those dependencies
RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

# Gopkg.toml and Gopkg.lock lists project dependencies
# These layers will only be re-built when Gopkg files are updated
COPY Gopkg.lock Gopkg.toml /go/src/github.com/klaper_/mqtt_data_exporter/
WORKDIR /go/src/github.com/klaper_/mqtt_data_exporter
# Install library dependencies
RUN dep ensure -vendor-only

# Copy all project and build it
# This layer will be rebuilt when ever a file has changed in the project directory
COPY . /go/src/github.com/klaper_/mqtt_data_exporter/
RUN go build -o /bin/mqtt_data_exporter

# This results in a single layer image
FROM alpine
COPY --from=build /bin/mqtt_data_exporter /bin/mqtt_data_exporter
EXPOSE 2112
ENTRYPOINT ["/bin/mqtt_data_exporter"]