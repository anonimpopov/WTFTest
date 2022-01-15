up:
	docker compose -f docker-compose.yml up --build --remove-orphans --detach

down:
	docker compose -f docker-compose.yml down --remove-orphans

log:
	docker logs wtftest-app-1
