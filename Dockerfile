FROM golang:latest

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab
# syntax=docker/dockerfile:1
#FROM golang:1.16-alpine AS build
#FROM larjim/kademlialab
WORKDIR /app

COPY d7024e/go.mod d7024e/go.sum ./
RUN go mod download

COPY /d7024e/*.go ./

RUN go build -o d7024e/main.go
RUN go build -o /D7024E-LAB-GROUP-12

CMD [ "/D7024E-LAB-GROUP-12" ]