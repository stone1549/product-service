package service

import (
	"context"
	"errors"
	"strings"

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
		first, err := strconv.Atoi(r.URL.Query().Get("first"))

		if err != nil {
			first = 20
			err = nil
		}

		cursor := r.URL.Query().Get("cursor")

		orderBy := common.OrderBy{}

		orderByStr := r.URL.Query().Get("orderBy")

		if orderByStr != "" {
			keyStrings := strings.Split(orderByStr, ",")

			for _, keyStr := range keyStrings {
				err = orderBy.Add(common.OrderByKey(keyStr))
				if err != nil {
					render.Render(w, r, errInvalidRequest(err))
					return
				}
			}
		}

		productRepo, ok := r.Context().Value("repo").(repository.ProductRepository)

		if !ok {
			render.Render(w, r, errRepository(errors.New("ProductRepository not found in context")))
			return
		}

		productsList, err := productRepo.GetProducts(r.Context(), first, cursor, orderBy)

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
