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
	cache2 "myservice/pkg/cache"
	smtp2 "myservice/pkg/smtp"
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

	cache := cache2.NewMemoryCache()
	smtp := smtp2.NewSMTP("smtp.gmail.com", "587", "rfsu vjub fnmm ditg", "golang.tester1974@gmail.com")

	userRepo := repository.NewUserRepository(pool)
	orderRepo := repository.NewOrderRepository(pool) 

	userSvc := service.NewUserService(userRepo, cache, smtp)
	orderSvc := service.NewOrderService(orderRepo)   

	userHandler := handlers.NewUserHandler(userSvc, logSvc)
	orderHandler := handlers.NewOrderHandler(orderSvc, logSvc) 

	

	mux := http.NewServeMux()
	
	mux.HandleFunc("POST /auth/register", userHandler.Register)
	mux.HandleFunc("POST /auth/login", userHandler.Login)
	
	// mux.Handle("POST /auth/login", middlware.Auth(http.HandleFunc()))
	mux.Handle("PUT /users/change-password", middlware.Auth(http.HandlerFunc(userHandler.ChangePassword)))
	mux.Handle("GET /users/profile", middlware.Auth(http.HandlerFunc(userHandler.GetMe)))
	mux.Handle("PATCH /orders/{id}/cancel", middlware.Auth(http.HandlerFunc(orderHandler.CancelOrder)))
	mux.Handle("GET /orders/my", middlware.Auth(http.HandlerFunc(orderHandler.GetMyOrders)))
	mux.Handle("POST /orders", middlware.Auth(http.HandlerFunc(orderHandler.Create)))
	mux.Handle("POST /verify", http.HandlerFunc(userHandler.Verify))
	mux.Handle("POST /auth/refresh", http.HandlerFunc(userHandler.Refresh))
	mux.Handle("POST /auth/logout", http.HandlerFunc(userHandler.Logout))
	
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

	logSvc.Info("Application successfully started and listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logSvc.Error("HTTP Server stopped: %v", err)
	}
}