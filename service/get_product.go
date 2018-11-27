package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stone1549/product-service/common"
	"github.com/stone1549/product-service/repository"
	"net/http"
)

type productResponse struct {
	Id               string     `json:"id"`
	Name             string     `json:"name"`
	DisplayImage     *string    `json:"displayImage"`
	Thumbnail        *string    `json:"thumbnail"`
	Price            *string    `json:"price"`
	Description      *string    `json:"description"`
	ShortDescription *string    `json:"shortDescription"`
	Quantity         int        `json:"quantity"`
	CreatedAt        *time.Time `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

func (plr productResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newProductResponse(product common.Product) productResponse {
	var price *string

	if product.Price == nil {
		str := product.Price.StringFixed(2)
		price = &str
	}

	return productResponse{
		Description:      product.Description,
		Name:             product.Name,
		DisplayImage:     product.DisplayImage,
		Id:               product.Id,
		Price:            price,
		Quantity:         product.QtyInStock,
		ShortDescription: product.ShortDescription,
		Thumbnail:        product.Thumbnail,
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
	}
}

func GetProductMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "productId")

		if id == "" {
			render.Render(w, r, errNotFound)
		}

		productRepo, ok := r.Context().Value("repo").(repository.ProductRepository)

		if !ok {
			render.Render(w, r, errRepository(errors.New("ProductRepository not found in context")))
			return
		}

		product, err := productRepo.GetProduct(r.Context(), id)

		if err != nil {
			render.Render(w, r, errRepository(err))
		} else if product == nil {
			render.Render(w, r, errNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "product", *product)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	product, ok := ctx.Value("product").(common.Product)

	if !ok {
		render.Render(w, r, errUnknown(errors.New("unable to retrieve product at this time")))
		return
	}

	if err := render.Render(w, r, newProductResponse(product)); err != nil {
		render.Render(w, r, errUnknown(err))
		return
	}
}
