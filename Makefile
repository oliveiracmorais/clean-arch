run:
	cd cmd/ordersystem \
	&& go run main.go wire_gen.go \
	&& cd ../..

migrate:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up
