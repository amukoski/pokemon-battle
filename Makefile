.PHONY: up down test logs clean

up:
	docker compose up --build -d

down:
	docker compose down

test:
	cd backend && go test ./...

logs:
	docker compose logs -f

clean:
	docker compose down -v
