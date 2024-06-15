package main

import (
	"log"
	"net/http"
	"time"

	"github.com/saadi925/flow"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	app := flow.NewRouter()
	register := flow.NewDependencyRegistry()
	register.Provide("user", "DatabaseUser")
	app.Use(register.Middleware)

	admin := app.Group("/admin")
	user := app.Group("/user")
	user.GET("/user", handleUserByID)
	admin.GET("/users/:id", handleUserByID)
	admin.GET("/dashboard", handleDashboard)
	config := flow.ServerConfig{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	server := flow.NewServer(app, config)
	go flow.GracefulShutdown(server)

	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func handleDashboard(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log.Println(ctx)
	w.Write([]byte("Handling dashboard..."))
}
func handleUserByID(w http.ResponseWriter, req *http.Request) {
	params := flow.GetParams(req)
	id := flow.ToInt(params["id"])
	dependency, err := flow.GetDependencyFromContext(req.Context(), "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	myUser, ok := dependency.(User)
	if !ok {
		flow.JSONError(w, "Error in Conversion , expected a user from user dependency", http.StatusBadRequest)
		return
	}
	queryParams := flow.ParseQueryParams(req)
	name := queryParams.Get("name")
	email := queryParams.Get("email")
	age := queryParams.GetInt("age")
	user := User{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}

	flow.JSONResponse(w, map[string]interface{}{
		"user":       user,
		"dependency": myUser,
	}, http.StatusCreated)
}
