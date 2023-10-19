# D7024E-Lab-Group-12

## Installation

[Install docker & docker-compose](https://docs.docker.com/)

[Install make](https://www.gnu.org/software/make/)
(Comes pre-installed with most Linux distros, `make -v`)

[Install Golang](https://golang.org/doc/install)

## Build & Run

This repo uses make to build and run the project.

`make build`: Builds the Docker image

`make run`: Runs the Docker image on a swarm

`make start`: Runs the `make build` and `make run` commands

`make stop`: Stops the swarm.

`make reload`: Runs `make stop` -> `make build` -> `make run`

`make test`: Runs the tests

For info on which commands to run, see the Makefile.

## UML Diagram

![UML Diagram](resource/d7024uml.svg)

(Last updated 2023-10-19)
