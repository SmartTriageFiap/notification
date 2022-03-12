MG_USER="admin"
MG_PASS="admin"
MG_ADDR="localhost"
MG_PORT="27017"

go-test:
	go test ./... -coverprofile cover.out

go-test-cover:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out

go-build:
	go build -o bin/main main.go

go-run:
	go run main.go

docker-stop:
	docker stop $$(docker ps -a -q)
	docker rm smart_triage_mongo
	docker rm zookeeper
	docker rm kafka

docker-kafka:
	docker run -d --net=host --name=zookeeper -e ZOOKEEPER_CLIENT_PORT=2181 confluentinc/cp-zookeeper:5.0.0
	docker run -d --net=host --name=kafka -e KAFKA_ZOOKEEPER_CONNECT=localhost:2181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 confluentinc/cp-kafka:5.0.0

docker-mongo:
	docker run --name smart_triage_mongo -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=${MG_USER} -e MONGO_INITDB_ROOT_PASSWORD=${MG_PASS} -d mongo:latest

docker-mongo-cli:
	docker run -it --rm --network host mongo mongo --host localhost -u admin -p admin --authenticationDatabase admin example

docker-build:
	docker build -t notification-app .

docker-exec:
	docker run --network host -p 8080 -e MONGO_ADDR=${MG_ADDR} -e MONGO_PORT=${MG_PORT} -e MONGO_USER=${MG_USER} -e MONGO_PASS=${MG_PASS}  geolocation-app:latest

go-exec:
	./bin/main

compile:
	GOOS=windows go build -o bin/main-windows main.go
	GOOS=linux go build -o bin/main-linux main.go

all: go-test go-build go-exec

all-docker: docker-stop docker-mongo docker-build docker-exec