SCRIPT := scripts/generate_env.py

all: up

build:
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.postgres.yml build

up: postgres

postgres:
	python ${SCRIPT}
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.postgres.yml up
	# cat .env | grep -v '^$$' | sed 's/^\(\w\)/export \1/g' | sed 's/$$/;/g' | xargs && docker compose

mongo:
	python ${SCRIPT}
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.mongo.yml up

down:
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.postgres.yml down
