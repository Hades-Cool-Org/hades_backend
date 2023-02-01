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
const dateStartParam = "date_start"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Create)                                      // assignar um pedido a um entrador
		r.Get("/", u.GetCustom)                                    // assignar um pedido a um entrador
		r.Delete("/{delivery_id}", u.Delete)                       // mudar o estado do pedido para coletado
		r.Get("/{delivery_id}", u.Get)                             // mudar o estado do pedido para coletado
		r.Post("/{delivery_id}/complete", u.Complete)              // recebebimento do pedido pelo gerente da loja
		r.Post("/{delivery_id}/user/{user_id}/collect", u.Collect) // mudar o estado do pedido para coletado

		r.Post("/user/{user_id}/start", u.StartDelivery) //Associar um carro a um entregador
		r.Post("/user/{user_id}/end", u.EndDelivery)     //end user turn
		r.Get("/user/{user_id}", u.GetByUser)            //get all deliveries by user
	}
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {

	deliveryId := chi.URLParam(r, deliveryIdParam)

	if deliveryId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is empty")))
		return
	}

	deliveries := &Delivery{
		ID: "ID_RETORNADO_DO_BANCO",
		Order: &Order{
			ID:        "ID_RETORNADO_DO_BANCO",
			StartDate: time.Now().Format(time.RFC3339),
			User: &User{
				ID:    "ID_RETORNADO_DO_BANCO",
				Name:  "user1",
				Email: "",
				Phone: "",
			},
			Products: []*Product{
				{
					ID:            "ID_RETORNADO_DO_BANCO",
					Name:          "product1",
					Image:         "image",
					MeasuringUnit: "cx",
					Quantity:      1,
					Total:         "10",
				},
				{
					ID:            "ID_RETORNADO_DO_BANCO",
					Name:          "product2",
					Image:         "url",
					MeasuringUnit: "cx",
					Quantity:      1,
					Total:         "",
				},
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{deliveries})
}

func (u *Router) GetCustom(w http.ResponseWriter, r *http.Request) {

	deliveryId := chi.URLParam(r, deliveryIdParam)

	if deliveryId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is empty")))
		return
	}

	dateStart := chi.URLParam(r, dateStartParam)
	if dateStart == "" {
		dateStart = time.Now().Format(time.RFC3339)
	}

	deliveries := []*Delivery{
		{
			ID: "ID_RETORNADO_DO_BANCO",
			Order: &Order{
				ID:        "ID_RETORNADO_DO_BANCO",
				StartDate: time.Now().Format(time.RFC3339),
				User: &User{
					ID:    "ID_RETORNADO_DO_BANCO",
					Name:  "user1",
					Email: "",
					Phone: "",
				},
				Products: []*Product{
					{
						ID:            "ID_RETORNADO_DO_BANCO",
						Name:          "product1",
						Image:         "image",
						MeasuringUnit: "cx",
						Quantity:      1,
						Total:         "10",
					},
					{
						ID:            "ID_RETORNADO_DO_BANCO",
						Name:          "product2",
						Image:         "url",
						MeasuringUnit: "cx",
						Quantity:      1,
						Total:         "",
					},
				},
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ListResponse{deliveries})
}

func (u *Router) GetByUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, userIdParam)

	if userId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("userId is empty")))
		return
	}

	deliveries := []*Delivery{
		{
			ID: "ID_RETORNADO_DO_BANCO",
			Order: &Order{
				ID:        "ID_RETORNADO_DO_BANCO",
				StartDate: time.Now().Format(time.RFC3339),
				User: &User{
					ID:    userId,
					Name:  "user1",
					Email: "user@gmail.com",
					Phone: "",
				},
				Products: []*Product{
					{
						ID:   "ID_RETORNADO_DO_BANCO",
						Name: "product1",
					},
				},
			},
		},
		{
			ID: "ID_RETORNADO_DO_BANCO2",
			Order: &Order{
				ID:        "ID_RETORNADO_DO_BANCO2",
				StartDate: time.Now().Format(time.RFC3339),
				User: &User{
					ID:    userId,
					Name:  "user1",
					Email: "user@gmail.com",
					Phone: "",
				},
				Products: []*Product{
					{
						ID:   "ID_RETORNADO_DO_BANCO",
						Name: "product1",
					},
				},
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ListResponse{deliveries})
}

func (u *Router) StartDelivery(w http.ResponseWriter, r *http.Request) {
	data := &StartDeliveryRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	userId := chi.URLParam(r, userIdParam)

	if userId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("userId is empty")))
		return
	}

	//db save and store

	render.Status(r, http.StatusOK)
}

func (u *Router) EndDelivery(w http.ResponseWriter, r *http.Request) {
	data := &StartDeliveryRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	//db save and store

	render.Status(r, http.StatusOK)
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

	userId := chi.URLParam(r, userIdParam)

	if userId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("userId is empty")))
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
