package controllers

import "github.com/duckcbuzz/crudapi/api/middlewares"

func (server *Server) initializeRoutes() {
	server.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(server.Home)).Methods("GET")

	server.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(server.CreateUser)).Methods("POST")
}
