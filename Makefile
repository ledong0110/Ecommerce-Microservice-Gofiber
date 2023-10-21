database_up:
	docker compose up db -d

database_down:
	docker compose stop db

auth_service_up:
	docker compose up auth_service -d

auth_service_down:
	docker compose stop auth_service

all_service_up:
	docker compose up -d

all_service_down:
	docker compose down



