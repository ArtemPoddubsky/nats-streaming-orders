all:
	@echo "Building..."
	@mkdir -p /var/tmp/docker/L0/postgresql
	@docker compose --env-file ./.env build
	@docker compose --env-file ./.env up -d

run:
	@docker compose --env-file ./.env up -d

publish:
	@$(MAKE) -C nats-service/ publish

test:
	@$(MAKE) -C nats-service/ test

app:
	@docker compose --env-file ./.env up -d app

nats:
	@docker compose --env-file ./.env up -d nats-streaming

db:
	@docker compose --env-file ./.env up -d database

clean:
	docker compose down
	@-docker volume rm $$(docker volume ls -q)
	@-docker rmi $$(docker images -q)

fclean:
	@-rm -rf /var/tmp/docker/L0/postgresql

re: clean all