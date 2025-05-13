

migrate-up:
	migrate -database "postgresql://root:root@localhost/root?sslmode=disable" -path migrations up

migrate-down:
	migrate -database "postgresql://root:root@localhost/root?sslmode=disable" -path migrations down
