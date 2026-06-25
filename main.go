package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"myservice/handlers"
	"myservice/internal/repository"
	"myservice/internal/service"
	"myservice/middlware"
	"myservice/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logSvc := logger.NewLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		"localhost", 5432, "postgres", "1234", "mydb"))
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	userRepo := repository.NewUserRepository(pool)
	orderRepo := repository.NewOrderRepository(pool) 

	userSvc := service.NewUserService(userRepo)
	orderSvc := service.NewOrderService(orderRepo)   

	userHandler := handlers.NewUserHandler(userSvc, logSvc)
	orderHandler := handlers.NewOrderHandler(orderSvc, logSvc) 

	mux := http.NewServeMux()
	
	mux.HandleFunc("/auth/register", userHandler.Register)
	mux.HandleFunc("/auth/login", userHandler.Login)

	mux.Handle("/users/me", middlware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetMe(w, r)
		case http.MethodPut:
			userHandler.UpdateMe(w, r)
		case http.MethodDelete:
			userHandler.DeleteMe(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/orders", middlware.Auth(http.HandlerFunc(orderHandler.Create)))

	logSvc.Info("Application successfully started and listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logSvc.Error("HTTP Server stopped: %v", err)
	}
}