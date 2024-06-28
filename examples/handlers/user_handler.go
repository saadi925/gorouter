package handlers

import (
	"fmt"
	"net/http"

	"github.com/saadi925/gorouter"
	"github.com/saadi925/gorouter/examples/dependencies"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func UserByID(w http.ResponseWriter, req *http.Request) {
	var role string = "user"
	params := gorouter.GetParams(req)
	id := params["id"]
	service1, err := gorouter.GetDependency(req.Context(), "globalServiceKey")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expectDB, err := gorouter.GetDependency(req.Context(), "db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// lets check if it is a db
	db, ok := expectDB.(dependencies.DB)
	if !ok {
		gorouter.JSONResponse(w, db, http.StatusInternalServerError)
		return
	}
	private_admin_service, err := gorouter.GetDependency(req.Context(), "private_admin_service")
	if err != nil {

		fmt.Println("admin_secret_not found :", err)
	}
	private_user_service, err := gorouter.GetDependency(req.Context(), "private_user_service")
	if err != nil {
		role = "admin"
		fmt.Println("user_secret_not found :", err)
	}

	queryParams := gorouter.ParseQueryParams(req)
	name := queryParams.Get("name")
	email := queryParams.Get("email")
	age := queryParams.GetInt("age")
	user := User{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}

	gorouter.JSONResponse(w, map[string]interface{}{
		"user":                  user,
		"globalService1":        service1,
		"db":                    db,
		"role":                  role,
		"private_admin_service": private_admin_service,
		"private_user_service":  private_user_service,
	}, http.StatusCreated)
}

func HandleDashboard(w http.ResponseWriter, req *http.Request) {
	private_admin_service, err := gorouter.GetDependency(req.Context(), "private_admin_service")
	if err != nil {
		fmt.Println("admin_secret_not found :", err)
	}
	key := gorouter.ParseQueryParams(req).Get("access_key")
	if key != private_admin_service {
		gorouter.JSONResponse(w, "Access Denied", http.StatusUnauthorized)
		return
	}
	gorouter.JSONResponse(w, "Access Granted", http.StatusAccepted)

}
