package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// Create a variable to hold my middleware chain
	standardMiddleware := alice.New(
		app.recoverPanicMiddleware,
		app.logRequestMiddleware,
		securityHeadersMiddleware
	)
	dynamicMiddleware := alice.New(app.session.Enable)
	// A third-party router/multiplexer
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/item/create", dynamicMiddleware.ThenFunc(app.createItem))
	mux.Post("/item/create", dynamicMiddleware.ThenFunc(app.createItem))
	mux.Get("/item/:id", dynamicMiddleware.ThenFunc(app.showItem))
	// Create a fileserver to serve our static content
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
