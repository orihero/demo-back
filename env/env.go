package env

import "github.com/gorilla/mux"

//--------GLOBAL VARIABLES---------------

var (
	Router    *mux.Router
	SecretKey = "secretkeyjwt"
)
