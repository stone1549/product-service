package service

import (
	"context"
	"errors"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stone1549/product-service/common"
	"github.com/stone1549/product-service/repository"
	"net/http"
	"strconv"
)

func SearchProductsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		first, err := strconv.Atoi(chi.URLParam(r, "first"))

		if err != nil {
			first = 20
		}

		cursor := r.URL.Query().Get("cursor")

		searchTxt := r.URL.Query().Get("searchTxt")

		productRepo, err := repository.GetProductRepository()

		if err != nil {
			render.Render(w, r, ErrRepository(err))
			return
		}

		productsList, err := productRepo.SearchProducts(r.Context(), searchTxt, first, cursor)

		if err != nil {
			render.Render(w, r, ErrRepository(err))
			return
		} else if len(productsList.Products) == 0 {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "products", productsList.Products)
		ctx = context.WithValue(ctx, "cursor", productsList.Cursor)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, ok := ctx.Value("products").([]common.Product)

	if !ok {
		render.Render(w, r, ErrUnknown(errors.New("unable to retrieve products at this time")))
		return
	}

	cursor := r.Context().Value("cursor").(string)

	if err := render.Render(w, r, newProductListResponse(products, cursor)); err != nil {
		render.Render(w, r, ErrUnknown(err))
		return
	}
}