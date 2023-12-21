
migration-up:
	migrate -path migrations/postgres/ -database "postgresql://ilyosbektemirov:123-4@localhost:5432/clinic?sslmode=disable" -verbose up

gen-swag:
	swag init -g api/api.go -o api/docs

run:
	go run cmd/main.go