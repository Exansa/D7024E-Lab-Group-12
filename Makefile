# Simple build scripts

build:
	docker build . -t kadlab

run:
	docker swarm init
	docker stack deploy -c docker-compose.yml kadswarm

detach:
	docker stack rm kadswarm
	docker swarm leave --force

reload:
	make detach
	make build
	make run
