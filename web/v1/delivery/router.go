package delivery

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/web/utils/net"
	"net/http"
	"time"
)

type Router struct {
}

func (u *Router) URL() string {
	return "/delivery"
}

const orderIdParam = "order_id"
const deliveryIdParam = "delivery_id"
const userIdParam = "user_id"
const productIdParam = "product_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Create)                         // assignar um pedido a um entrador
		r.Post("/{delivery_id}/collect", u.Collect)   // mudar o estado do pedido para coletado
		r.Delete("/{delivery_id}", u.Delete)          // mudar o estado do pedido para coletado
		r.Post("/{delivery_id}/complete", u.Complete) // recebebimento do pedido pelo gerente da loja

		r.Post("/", u.Iniciar) //Associar um carro a um entregador
	}
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {
	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	delivery := &Delivery{
		ID:    "id from db",
		State: "ABERTO",
		Order: &Order{
			ID:        data.Order.ID,
			StartDate: time.Now().Format(time.RFC3339),
			User: &User{
				ID:    "id",
				Name:  "comprador",
				Email: "comprador@gmail.com",
				Phone: "+551999999999",
			},
			Products: []*Product{
				{
					ID:            "from-db",
					Name:          "Tomate",
					Image:         "url",
					MeasuringUnit: "UN",
					Quantity:      10,
					Total:         "10.0",
				},
				{
					ID:            "from-db",
					Name:          "Tomate cereja",
					Image:         "url",
					MeasuringUnit: "UN",
					Quantity:      10,
					Total:         "59.2",
				},
			},
		},
		Store: &Store{
			ID:      "123",
			Name:    "Loja1",
			Address: "Rua 30 de julho",
		},
		User: &User{
			ID:    data.User.ID,
			Name:  "entregador",
			Email: "entregador@gmail.com",
			Phone: "+551999999999",
		},
		Vehicle: &Vehicle{
			ID:   "id",
			Name: "fiorino",
		},
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{delivery})
}

func (u *Router) Complete(w http.ResponseWriter, r *http.Request) {
	data := &CompleteDeliveryRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	completeDelivery := &CompleteDeliveryRequest{
		ID:       "from db",
		Products: data.Products,
	}

	//do db magic

	render.Status(r, http.StatusOK)
	render.Render(w, r, &CompleteDeliveryResponse{completeDelivery})
}

func (u *Router) Collect(w http.ResponseWriter, r *http.Request) {

	deliveryId := chi.URLParam(r, deliveryIdParam)

	if deliveryId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is empty")))
		return
	}

	//change status no db

	render.Status(r, http.StatusOK)
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

	deliveryId := chi.URLParam(r, deliveryIdParam)

	if deliveryId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is empty")))
		return
	}

	//delete no db

	render.Status(r, http.StatusNoContent)
}
