package main

import (
	"context"
	"net/http"
	"time"

	"myservice/handlers"
	"myservice/internal/repository"
	"myservice/internal/service"
	"myservice/pkg/db"
	"myservice/pkg/logger"
)

func main(){
	logSvc := logger.NewLogger()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := "postgres://postgres:pass@localhost:5432/mydb"
	pool, err := db.InitDB(ctx, dsn)
	if err != nil{
		logSvc.Error("DB connection pool initialization failed: %v", err)
		return
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo)
	UserHandler := handlers.NewUserHandler(userSvc, logSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("auth/register", UserHandler.Register)

	logSvc.Info("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil{
		logSvc.Error("HTTP server error: %v", err)
	}
}