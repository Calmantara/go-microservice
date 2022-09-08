infra:
	cd docker-infra && docker-compose -f docker-compose.yaml up -d
up:
	echo "build bin files"
	cd go-wallet/cmd && env ENV=PROD GOOS=linux GOARCH=amd64 go build -o ../main-app
	cd go-emitter/cmd && env ENV=PROD GOOS=linux GOARCH=amd64 go build -o ../main-app
	cd go-gateway/cmd &&  env ENV=PROD GOOS=linux GOARCH=amd64 go build -o ../main-app
	echo "running docker compose"
	docker-compose -f docker-compose.yaml up -d
down:
	rm go-wallet/main-app
	rm go-emitter/main-app
	rm go-gateway/main-app
	cd docker-infra && docker-compose -f docker-compose.yaml down
	docker-compose -f docker-compose.yaml down
topic:
	docker exec -it broker kafka-topics --create --topic balance-transaction --bootstrap-server localhost:9092 
	docker exec -it broker kafka-topics --create --topic balance-transaction-group-table --bootstrap-server localhost:9092