package order

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/app/web/utils/net"
	"net/http"
	"time"
)

type Router struct {
}

func (u *Router) URL() string {
	return "/orders"
}

const orderIdParam = "order_id"
const userIdParam = "user_id"
const productIdParam = "product_id"

// const storeIdParam = "store_id"
const stateParam = "state"
const dateStartParam = "dateStart"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Create) //todo: should we update?
		r.Get("/", u.GetCustom)
		r.Delete("/{order_id}", u.Delete)                             //todo: check permissions
		r.Post("/{order_id}/product", u.AddProduct)                   //add product // TODO: verificar se aceitar multiplos é uma boa ideia
		r.Put("/{order_id}/product/{product_id}", u.UpdateProduct)    //add product
		r.Get("/{order_id}/product/{product_id}", u.GetProduct)       //add product
		r.Delete("/{order_id}/product/{product_id}", u.DeleteProduct) //add product
	}
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	//db save
	order := &Order{
		ID: "id_from_db",
		Vendor: &Vendor{
			ID:      "from_db",
			Name:    "Joao do Tomate",
			Address: "Rua 30 de Julho",
		},
		StartDate: time.Now().Format(time.RFC3339),
		State:     "CREATED",
		User: &User{
			ID:   "from-db",
			Name: "Oscar",
		},
		Total: "69.2",
		Payments: []*Payment{
			{
				Type:  "MONEY",
				Total: "10.0",
				Date:  time.Now().Format(time.RFC3339),
			},
			{
				Type:  "PIX",
				Total: "59.2",
				Date:  time.Now().Format(time.RFC3339),
			},
		},
		Products: []*Product{
			{
				ID:            "from-db",
				Name:          "Tomate",
				Image:         "url",
				MeasuringUnit: "UN",
				Quantity:      10,
				Total:         "10.0",
				Stores: []*Store{
					{
						ID:       "from-db",
						Name:     "Loja1",
						Address:  "Rua 30 de julho",
						Quantity: 10,
					},
				},
			},
			{
				ID:            "from-db",
				Name:          "Tomate cereja",
				Image:         "url",
				MeasuringUnit: "UN",
				Quantity:      10,
				Total:         "59.2",
				Stores: []*Store{
					{
						ID:       "from-db",
						Name:     "Loja1",
						Address:  "Rua 30 de julho",
						Quantity: 10,
					},
				},
			},
		},
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{order})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

	orderId := chi.URLParam(r, orderIdParam)

	if orderId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	//db delete
	render.Status(r, http.StatusNoContent)
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {

	orderId := chi.URLParam(r, orderIdParam)

	if orderId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	//db get
	order := &Order{
		ID: "id_from_db",
		Vendor: &Vendor{
			ID:      "from_db",
			Name:    "Joao do Tomate",
			Address: "Rua 30 de Julho",
		},
		StartDate: time.Now().Format(time.RFC3339),
		State:     "CREATED",
		User: &User{
			ID:   "from-db",
			Name: "Oscar",
		},
		Total: "69.2",
		Payments: []*Payment{
			{
				Type:  "MONEY",
				Total: "10.0",
				Date:  time.Now().Format(time.RFC3339),
			},
			{
				Type:  "PIX",
				Total: "59.2",
				Date:  time.Now().Format(time.RFC3339),
			},
		},
		Products: []*Product{
			{
				ID:            "from-db",
				Name:          "Tomate",
				Image:         "url",
				MeasuringUnit: "UN",
				Quantity:      10,
				Total:         "10.0",
				Stores: []*Store{
					{
						ID:       "from-db",
						Name:     "Loja1",
						Address:  "Rua 30 de julho",
						Quantity: 10,
					},
				},
			},
			{
				ID:            "from-db",
				Name:          "Tomate cereja",
				Image:         "url",
				MeasuringUnit: "UN",
				Quantity:      10,
				Total:         "59.2",
				Stores: []*Store{
					{
						ID:       "from-db",
						Name:     "Loja1",
						Address:  "Rua 30 de julho",
						Quantity: 10,
					},
				},
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{order})
}

func (u *Router) GetCustom(w http.ResponseWriter, r *http.Request) {

	state := chi.URLParam(r, stateParam)

	dateStart := chi.URLParam(r, dateStartParam)
	if dateStart == "" {
		dateStart = time.Now().Format(time.RFC3339)
	}

	userId := chi.URLParam(r, userIdParam)

	//db get
	order := &Order{
		ID: "id_from_db",
		Vendor: &Vendor{
			ID:      "from_db",
			Name:    "Joao do Tomate",
			Address: "Rua 30 de Julho",
		},
		StartDate: dateStart,
		State:     state,
		User: &User{
			ID:   userId,
			Name: "Oscar",
		},
		Total: "69.2",
		Payments: []*Payment{
			{
				Type:  "MONEY",
				Total: "10.0",
				Date:  time.Now().Format(time.RFC3339),
			},
			{
				Type:  "PIX",
				Total: "59.2",
				Date:  time.Now().Format(time.RFC3339),
			},
		},
		Products: []*Product{
			{
				ID:            "from-db",
				Name:          "Tomate",
				Image:         "url",
				MeasuringUnit: "UN",
				Quantity:      10,
				Total:         "10.0",
				Stores: []*Store{
					{
						ID:       "from-db",
						Name:     "Loja1",
						Address:  "Rua 30 de julho",
						Quantity: 10,
					},
				},
			},
			{
				ID:            "from-db",
				Name:          "Tomate cereja",
				Image:         "url",
				MeasuringUnit: "UN",
				Quantity:      10,
				Total:         "59.2",
				Stores: []*Store{
					{
						ID:       "from-db",
						Name:     "Loja1",
						Address:  "Rua 30 de julho",
						Quantity: 10,
					},
				},
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ListResponse{[]*Order{order}})
}

func (u *Router) AddProduct(w http.ResponseWriter, r *http.Request) {

	data := &AddProductRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	orderId := chi.URLParam(r, orderIdParam)

	if orderId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	//db add

	order := &Order{
		ID: "id_from_db",
		Vendor: &Vendor{
			ID:      "from_db",
			Name:    "Joao do Tomate",
			Address: "Rua 30 de Julho",
		},
		StartDate: time.Now().Format(time.RFC3339),
		State:     "CREATED",
		User: &User{
			ID:   "from-db",
			Name: "Oscar",
		},
		Total: "69.2",
		Payments: []*Payment{
			{
				Type:  "MONEY",
				Total: "10.0",
				Date:  time.Now().Format(time.RFC3339),
			},
			{
				Type:  "PIX",
				Total: "59.2",
				Date:  time.Now().Format(time.RFC3339),
			},
		},
		Products: []*Product{
			{
				ID:            "from-db",
				Name:          "Tomate",
				Image:         "url",
				MeasuringUnit: "UN",
				Quantity:      10,
				Total:         "10.0",
				Stores: []*Store{
					{
						ID:       "from-db",
						Name:     "Loja1",
						Address:  "Rua 30 de julho",
						Quantity: 10,
					},
				},
			},
			{
				ID:            "from-db",
				Name:          "Tomate cereja",
				Image:         "url",
				MeasuringUnit: "UN",
				Quantity:      10,
				Total:         "59.2",
				Stores: []*Store{
					{
						ID:       "from-db",
						Name:     "Loja1",
						Address:  "Rua 30 de julho",
						Quantity: 10,
					},
				},
			},
		},
	}

	for _, product := range data.Products {
		order.Products = append(order.Products, product)
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{order})
}

func (u *Router) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	data := &UpdateProductRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	orderId := chi.URLParam(r, orderIdParam)

	if orderId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	data.ID = productId

	order := &Order{
		ID: "id_from_db",
		Vendor: &Vendor{
			ID:      "from_db",
			Name:    "Joao do Tomate",
			Address: "Rua 30 de Julho",
		},
		StartDate: time.Now().Format(time.RFC3339),
		State:     "CREATED",
		User: &User{
			ID:   "from-db",
			Name: "Oscar",
		},
		Total: "69.2",
		Payments: []*Payment{
			{
				Type:  "MONEY",
				Total: "10.0",
				Date:  time.Now().Format(time.RFC3339),
			},
			{
				Type:  "PIX",
				Total: "59.2",
				Date:  time.Now().Format(time.RFC3339),
			},
		},
		Products: []*Product{
			{
				ID:            "from-db",
				Name:          "Tomate",
				Image:         "url",
				MeasuringUnit: "UN",
				Quantity:      10,
				Total:         "10.0",
				Stores: []*Store{
					{
						ID:       "from-db",
						Name:     "Loja1",
						Address:  "Rua 30 de julho",
						Quantity: 10,
					},
				},
			},
			{
				ID:            "from-db",
				Name:          "Tomate cereja",
				Image:         "url",
				MeasuringUnit: "UN",
				Quantity:      10,
				Total:         "59.2",
				Stores: []*Store{
					{
						ID:       "from-db",
						Name:     "Loja1",
						Address:  "Rua 30 de julho",
						Quantity: 10,
					},
				},
			},
		},
	}

	order.Products = append(order.Products, data.Product)

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{order})
}

func (u *Router) GetProduct(w http.ResponseWriter, r *http.Request) {

	orderId := chi.URLParam(r, orderIdParam)

	if orderId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	product := &Product{
		ID:            "from-db",
		Name:          "Tomate",
		Image:         "url",
		MeasuringUnit: "UN",
		Quantity:      10,
		Total:         "10.0",
		Stores: []*Store{
			{
				ID:       "from-db",
				Name:     "Loja1",
				Address:  "Rua 30 de julho",
				Quantity: 10,
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ProductResponse{product})
}

func (u *Router) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	orderId := chi.URLParam(r, orderIdParam)

	if orderId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}
	//db delete
	render.Status(r, http.StatusNoContent)
}
