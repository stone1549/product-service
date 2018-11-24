package service

import (
	"context"
	"errors"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stone1549/product-service/common"
	"github.com/stone1549/product-service/repository"
	"net/http"
)

type productResponse struct {
	Id               string  `json:"id"`
	Name             string  `json:"name"`
	DisplayImage     *string `json:"displayImage"`
	Thumbnail        *string `json:"thumbnail"`
	Price            *string `json:"price"`
	Description      *string `json:"description"`
	ShortDescription *string `json:"shortDescription"`
	Quantity         int     `json:"quantity"`
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
			render.Render(w, r, ErrRepository(err))
			return
		}

		product, err := productRepo.ProductFromRepo(r.Context(), id)

		if err != nil {
			render.Render(w, r, ErrRepository(err))
		} else if product == nil {
			render.Render(w, r, ErrNotFound)
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
		render.Render(w, r, ErrUnknown(errors.New("unable to retrieve product at this time")))
		return
	}

	if err := render.Render(w, r, newProductResponse(product)); err != nil {
		render.Render(w, r, ErrUnknown(err))
		return
	}
}
