build:
	@go build -o server.exe .


test:
	@go test -v ./...

run:
	@air

migrate-dev: # add migration name at the end (ex: make migration create-cars-table)
	@npx prisma migrate dev

migrate-up:
	@npx prisma migrate up

# Aliases for database-related commands
db\:dev: migrate-dev
db\:deploy: migrate-up



