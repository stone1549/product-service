package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/stone1549/product-service/repository"
	"github.com/stone1549/product-service/service"
	"net/http"
	"time"
)

func main() {
	flag.Parse()

	err := repository.ConfigureProductRepository(repository.InMemory)

	if err != nil {
		panic(fmt.Sprintf("Unable to configure repository: %s", err.Error()))
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/products", func(r chi.Router) {
		r.With(service.ListProductsMiddleware).Get("/", service.ListProducts)
		//r.Route("/{productId}", func(r chi.Router) {
		//	r.Use(product.ProductMiddleware) // Load the *Article on the request context
		//	r.Get("/", GetProduct)           // GET /articles/123
		//})
	})

	http.ListenAndServe(":3333", r)
}
