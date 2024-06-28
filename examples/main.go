package main

import (
	"log"
	"time"

	"github.com/saadi925/gorouter"
	"github.com/saadi925/gorouter/examples/dependencies"
	"github.com/saadi925/gorouter/examples/handlers"
)

const AdminService = "privateAdminService"

func main() {
	router := gorouter.NewRouter()
	// Creating A New Dependency Register
	register := gorouter.NewDependencyRegistry()
	// Register Global Dependencies
	register.Provide("globalServiceKey", "GlobalService1")
	register.Provide("db", dependencies.InitDB())
	// Using Register Middleware
	router.Use(register.Middleware)
	// Using Group (could be used for versioning)
	admin := router.Group("/admin")
	user := router.Group("/user")
	// Register private dependencies
	user.Provide("private_user_service", dependencies.UserService())
	admin.Provide("private_admin_service", dependencies.AdminService())

	admin.GET("/user/:id", handlers.UserByID)

	admin.GET("/dashboard", handlers.HandleDashboard)

	user.GET("/:id", handlers.UserByID)

	config := gorouter.ServerConfig{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	server := gorouter.NewServer(router, config)
	go gorouter.GracefulShutdown(server)

	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
