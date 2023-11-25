database_up:
	docker compose up db rabbitmq -d

database_down:
	docker compose stop db rabbitmq

auth_service_up:
	docker compose up auth_service -d

auth_service_down:
	docker compose stop auth_service

email_service_up:
	docker compose up email_service -d

email_service_down:
	docker compose stop email_service

nginx_up:
	docker compose up nginx -d

all_service_up:
	docker compose up -d

all_service_down:
	docker compose down



