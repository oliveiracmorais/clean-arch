package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/oliveiracmorais/clean-arch/configs"
	"github.com/oliveiracmorais/clean-arch/internal/event/handler"
	"github.com/oliveiracmorais/clean-arch/internal/infra/graph"

	"github.com/oliveiracmorais/clean-arch/internal/infra/grpc/pb"
	"github.com/oliveiracmorais/clean-arch/internal/infra/grpc/service"
	"github.com/oliveiracmorais/clean-arch/internal/infra/web/webserver"
	"github.com/oliveiracmorais/clean-arch/pkg/events"
	"github.com/streadway/amqp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Iniciando os servidores WEB, GraphQL, GRPC e RabbitMQ ...")

	config, err := configs.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(config.DBDriver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(config.AmpqURL, config.AmpqPort)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	eventDispatcher.Register("OrdersListed", &handler.OrdersListedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db, eventDispatcher)

	webserver := webserver.NewWebServer(":" + config.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/createOrder", webOrderHandler.Create)
	webserver.AddHandler("/order", webOrderHandler.List)
	fmt.Println("Starting web server on port", config.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", config.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", config.GraphQLServerPort)
	http.ListenAndServe(":"+config.GraphQLServerPort, nil)
}

func getRabbitMQChannel(url string, port string) *amqp.Channel {
	urlPort := fmt.Sprintf("%s:%s/", url, port)
	fmt.Println("url=", url, "port=", port)
	conn, err := amqp.Dial(urlPort)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
