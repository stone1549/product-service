package service

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stone1549/product-service/common"
	"github.com/stone1549/product-service/repository"
	"net/http"
)

type productResponse struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	DisplayImage     string `json:"displayImage"`
	Thumbnail        string `json:"thumbnail"`
	Price            string `json:"price"`
	Description      string `json:"description"`
	ShortDescription string `json:"shortDescription"`
	Quantity         int    `json:"quantity"`
}

func NewProductResponse(product *common.Product) productResponse {
	return productResponse{
		Description:      product.Description,
		Name:             product.Name,
		DisplayImage:     product.DisplayImage,
		Id:               product.Id,
		Price:            product.Price.StringFixed(2),
		Quantity:         product.Quantity,
		ShortDescription: product.ShortDescription,
		Thumbnail:        product.Thumbnail,
	}
}

func ProductMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "productId")

		if id == "" {
			render.Render(w, r, ErrNotFound)
		}

		productRepo, err := repository.GetProductRepository()

		if err != nil {
			render.Render(w, r, ErrUnknown(err))
		}

		product, err := productRepo.ProductFromRepo(id)

		if err != nil {
			render.Render(w, r, ErrUnknown(err))
		}

		ctx := context.WithValue(r.Context(), "product", product)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
