package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/saraqueirozs/clean-arch-challenge/configs"
	"github.com/saraqueirozs/clean-arch-challenge/internal/event/handler"
	"github.com/saraqueirozs/clean-arch-challenge/internal/infra/graph"
	"github.com/saraqueirozs/clean-arch-challenge/internal/infra/grpc/pb"
	"github.com/saraqueirozs/clean-arch-challenge/internal/infra/grpc/service"
	"github.com/saraqueirozs/clean-arch-challenge/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	// webserver := webserver.NewWebServer(configs.WebServerPort)
	// webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	// webserver.AddHandler("/order", webOrderHandler.Create)
	// fmt.Println("Starting web server on port", configs.WebServerPort)
	// go webserver.Start()
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	router := chi.NewRouter()
	router.Route("/order", func(r chi.Router) {
		r.Post("/", webOrderHandler.Create)
		r.Get("/", webOrderHandler.ListOrders)
	})
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go http.ListenAndServe(":"+configs.WebServerPort, router)

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
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

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", configs.RabbitMQUser, configs.RabbitMQPassword, configs.RabbitMQHost, configs.RabbitMQPort))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
