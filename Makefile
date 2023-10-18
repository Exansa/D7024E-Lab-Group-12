# Simple build scripts

build:
	docker build . -t kadlab

run:
	docker swarm init
	docker stack deploy -c docker-compose.yml kadswarm

stop:
	docker stack rm kadswarm
	docker swarm leave --force

reload:
	make stop
	make build
	make run

start:
	make build
	make run

test:
	$(MAKE) -C d7024e test
	