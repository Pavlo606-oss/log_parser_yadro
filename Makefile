DOCKER = docker-compose

.PHONY: up rebuild down down-clean logs processes migrate-up migrate-down

up:
	$(DOCKER) up -d

rebuild:
	$(DOCKER) up -d --build

down:
	$(DOCKER) down

down-clean:
	$(DOCKER) down -v

logs:
	$(DOCKER) logs

processes:
	$(DOCKER) ps -a

migrate-up:
	docker exec -i log_parser-postgres psql -U parser_user -d parser_db < migrations/0001_init.up.sql

migrate-down:
	docker exec -i log_parser-postgres psql -U parser_user -d parser_db < migrations/0001_init.down.sql