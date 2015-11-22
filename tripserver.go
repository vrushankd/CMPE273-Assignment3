package main

import (
	// Controllers are kept in this package
	"assgn3Controllers"
	// Standard library packages
	"net/http"
	// Third party packages
	"github.com/julienschmidt/httprouter"
)

func main() {
	// Instantiate a new router
	router := httprouter.New()

	// Get a controller instance
	controller := assgn3Controllers.NewUserController()

	// Add handlers
	router.GET("/trips/:id", controller.CheckTrip)
	router.POST("/trips", controller.PlanTrip)
	router.PUT("/trips/:id/request", controller.CheckNextDestination)

	// Expose the server at port 3000
	http.ListenAndServe(":3000", router)
}
