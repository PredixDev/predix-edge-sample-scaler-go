# start with the Predix Edge base alpine image
FROM golang

LABEL maintainer="Predix Edge Adoption"
LABEL hub="https://hub.docker.com"
LABEL org="https://hub.docker.com/u/predixedge"
LABEL repo="predix-edge-sample-scaler-go"
LABEL version="1.0.6"
LABEL support="https://forum.predix.io"
LABEL license="https://github.com/PredixDev/predix-docker-samples/blob/master/LICENSE.md"

# Create app directory in the image
WORKDIR /usr/src/predix-edge-sample-scaler-go

# copy app's source files to the image
COPY src/app.go .

RUN go get -u github.com/eclipse/paho.mqtt.golang
RUN ls
RUN go build app.go

#Start the app
CMD ["./app"]
