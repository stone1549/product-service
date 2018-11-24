package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/stone1549/product-service/common"
	"github.com/stone1549/product-service/repository"
	"github.com/stone1549/product-service/service"
	"net/http"
)

func main() {
	flag.Parse()

	config, err := common.GetConfiguration()

	if err != nil {
		panic(fmt.Sprintf("Unable to load configuration: %s", err.Error()))
	}

	err = repository.ConfigureProductRepository(config)

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
	r.Use(middleware.Timeout(config.GetTimeout()))

	r.Route("/products", func(r chi.Router) {
		r.With(service.ListProductsMiddleware).Get("/", service.ListProducts)
		r.Route("/{productId}", func(r chi.Router) {
			r.Use(service.ProductMiddleware) // Load the *Article on the request context
			r.Get("/", service.GetProduct)   // GET /articles/123
		})
	})

	http.ListenAndServe(":3333", r)
}
