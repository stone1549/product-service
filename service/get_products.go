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

type productListResponse struct {
	Products []productResponse `json:"products"`
	Cursor   string            `json:"cursor"`
}

func (plr productListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newProductListResponse(products []common.Product, cursor string) productListResponse {
	results := make([]productResponse, 0)
	for _, product := range products {
		productResponse := newProductResponse(product)
		results = append(results, productResponse)
	}

	return productListResponse{results, cursor}
}

func GetProductsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		first, err := strconv.Atoi(chi.URLParam(r, "first"))

		if err != nil {
			first = 20
		}

		cursor := r.URL.Query().Get("cursor")

		productRepo, err := repository.GetProductRepository()

		if err != nil {
			render.Render(w, r, errRepository(err))
			return
		}

		productsList, err := productRepo.GetProducts(r.Context(), first, cursor)

		if err != nil {
			render.Render(w, r, errRepository(err))
			return
		} else if len(productsList.Products) == 0 {
			render.Render(w, r, errNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "products", productsList.Products)
		ctx = context.WithValue(ctx, "cursor", productsList.Cursor)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, ok := ctx.Value("products").([]common.Product)

	if !ok {
		render.Render(w, r, errUnknown(errors.New("unable to retrieve products at this time")))
		return
	}

	cursor := r.Context().Value("cursor").(string)

	if err := render.Render(w, r, newProductListResponse(products, cursor)); err != nil {
		render.Render(w, r, errUnknown(err))
		return
	}
}
