package routes

import (
	"net/http"

	"../controllers"
	"../env"
	"../middleware"
	"github.com/gorilla/mux"
)

//----------------------ROUTES-------------------------------
//create a mux router
func CreateRouter() {
	env.Router = mux.NewRouter()
}

//initialize all routes
func InitializeRoute() {
	//Auth
	env.Router.HandleFunc("/signup", controllers.SignUp).Methods("POST")
	env.Router.HandleFunc("/signin", controllers.SignIn).Methods("POST")
	env.Router.HandleFunc("/user/block-unblock", middleware.IsAuthorized(controllers.BlockUnblockUser)).Methods("GET")
	env.Router.HandleFunc("/users", middleware.IsAuthorized(controllers.GetUsers)).Methods("GET")
	//Restaurant
	env.Router.HandleFunc("/restaurant", middleware.IsAuthorized(controllers.CreateRestaurant)).Methods("POST")
	env.Router.HandleFunc("/restaurant", middleware.IsAuthorized(controllers.GetAllRestaurants)).Methods("GET")
	env.Router.HandleFunc("/my-restaurants", middleware.IsAuthorized(controllers.GetUserRestaurants)).Methods("GET")
	env.Router.HandleFunc("/meal", middleware.IsAuthorized(controllers.CreateMeal)).Methods("POST")
	env.Router.HandleFunc("/restaurant-meals", middleware.IsAuthorized(controllers.GetRestaurantMeals)).Methods("GET")
	env.Router.HandleFunc("/change-order-status", middleware.IsAuthorized(controllers.ChangeOrderStatus)).Methods("GET")
	env.Router.HandleFunc("/user-orders", middleware.IsAuthorized(controllers.UserOrders)).Methods("GET")
	env.Router.HandleFunc("/restaurant-orders", middleware.IsAuthorized(controllers.RestaurantOrders)).Methods("GET")
	//Order
	env.Router.HandleFunc("/order", middleware.IsAuthorized(controllers.CreateOrder)).Methods("POST")
	//Other
	env.Router.HandleFunc("/user", middleware.IsAuthorized(controllers.UserIndex)).Methods("GET")
	env.Router.HandleFunc("/admin", middleware.IsAuthorized(controllers.AdminIndex)).Methods("GET")
	env.Router.HandleFunc("/upload", middleware.IsAuthorized(controllers.MultipleFileUpload)).Methods("POST")
	env.Router.HandleFunc("/download/{name}", controllers.GetUploadedFiles).Methods("GET")
	env.Router.HandleFunc("/", controllers.Index).Methods("GET")
	env.Router.PathPrefix("/public/uploads/").Handler(http.StripPrefix("/public/uplods/", http.FileServer((http.Dir("./public/uploads/")))))

}
