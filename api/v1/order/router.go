package order

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"hades_backend/api/utils/net"
	order2 "hades_backend/app/cmd/order"
	"hades_backend/app/model"
	"net/http"
	"strconv"
	"time"
)

type Router struct {
	db *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	err := db.AutoMigrate(&order2.Order{}, &order2.Item{}, &order2.Payment{})
	if err != nil {
		panic(err)
	}
	return &Router{db: db}
}

func (u *Router) URL() string {
	return "/orders"
}

const orderIdParam = "order_id"
const productIdParam = "product_id"
const paymentIdParam = "payment_id"

// const storeIdParam = "store_id"
const stateParam = "state"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Create) //todo: should we update?
		r.Get("/", u.GetAll)

		r.Get("/{order_id}", u.Get)       //todo: check permissions
		r.Delete("/{order_id}", u.Delete) //todo: check permissions
		r.Put("/{order_id}", u.Update)    //todo: check permissions

		r.Post("/{order_id}/payment", u.AddPayment)                   //todo: check permissions
		r.Delete("/{order_id}/payment/{payment_id}", u.RemovePayment) //todo: check permissions

		r.Get("/{order_id}/product/{product_id}", u.GetProduct)
		r.Delete("/{order_id}/product", u.DeleteProduct) //add product
	}
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {

	var request Request

	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	o, err := order2.CreateOrder(r.Context(), request.Order)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{convertOrderToResponse(o)})
}

func convertOrderToResponse(o *order2.Order) *model.Order {

	payments := make([]*model.Payment, len(o.Payments))

	for i, p := range o.Payments {
		payments[i] = &model.Payment{
			ID:    p.ID,
			Type:  p.Type,
			Total: p.Total,
			Date:  p.CreatedAt.Format(time.RFC3339),
			Text:  p.Text,
		}
	}

	is := convertOrderItems(o.Items)

	s, _ := model.OrderStateFromString(o.State)

	z := &model.Order{
		ID: o.ID,
		Vendor: &model.Vendor{
			ID:       o.Vendor.ID,
			Name:     o.Vendor.Name,
			Email:    o.Vendor.Email,
			Phone:    o.Vendor.Phone,
			Cnpj:     o.Vendor.Cnpj,
			Type:     o.Vendor.Type,
			Location: o.Vendor.Location,
			Contact: &model.Contact{
				Name:  o.Vendor.Contact.Name,
				Email: o.Vendor.Contact.Email,
				Phone: o.Vendor.Contact.Phone,
			},
		},
		CreatedDate: o.CreatedAt.Format(time.RFC3339),
		State:       &s,
		EndDate: func() *string {
			if o.CompletedDate != nil {
				o.CompletedDate.Format(time.RFC3339)
			}
			return nil
		}(),
		User: &model.User{
			ID:    o.User.ID,
			Name:  o.User.Name,
			Email: o.User.Email,
			Phone: o.User.Phone,
		},
		Total:    o.Total,
		Payments: payments,
		Items:    is,
	}

	return z
}

func convertOrderItems(z []*order2.Item) []*model.Item {
	is := make([]*model.Item, len(z))

	for i, p := range z {
		is[i] = &model.Item{
			ProductID:     p.ProductID,
			OrderID:       p.OrderID,
			StoreID:       p.StoreID,
			Name:          p.Product.Name,
			ImageUrl:      p.Product.ImageUrl,
			MeasuringUnit: p.Product.MeasuringUnit,
			Quantity:      p.Quantity,
			Available:     p.Available,
			Total:         p.CalculateTotal(),
		}
	}
	return is
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {

	chi.URLParam(r, stateParam)

	o := &order2.GetOrdersOptions{
		Params: r.URL.Query(),
	}

	orders, err := order2.GetOrders(r.Context(), o)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	respOrders := make([]*model.Order, len(orders))

	for i, o := range orders {
		respOrders[i] = convertOrderToResponse(o)
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ListResponse{respOrders})
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {

	oID := chi.URLParam(r, orderIdParam)

	if oID == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	oIDInt, err := strconv.Atoi(oID)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is not a number: "+err.Error())))
		return
	}

	o, err := order2.GetOrder(r.Context(), uint(oIDInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{convertOrderToResponse(o)})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

	oID := chi.URLParam(r, orderIdParam)

	if oID == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	oIDInt, err := strconv.Atoi(oID)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is not a number: "+err.Error())))
		return
	}

	err = order2.DeleteOrder(r.Context(), uint(oIDInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (u *Router) Update(w http.ResponseWriter, r *http.Request) {

	oID := chi.URLParam(r, orderIdParam)

	if oID == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	oIDInt, err := strconv.Atoi(oID)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is not a number: "+err.Error())))
		return
	}

	var request UpdateRequest

	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	request.ID = uint(oIDInt)

	err = order2.UpdateOrder(r.Context(), uint(oIDInt), request.Order)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}

func (u *Router) AddPayment(w http.ResponseWriter, r *http.Request) {

	oID := chi.URLParam(r, orderIdParam)

	if oID == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	oIDInt, err := strconv.Atoi(oID)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is not a number: "+err.Error())))
		return
	}

	var request PaymentRequest

	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	_, err = order2.AddPayment(r.Context(), uint(oIDInt), request.Payment)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}

func (u *Router) RemovePayment(w http.ResponseWriter, r *http.Request) {

	oID := chi.URLParam(r, orderIdParam)

	if oID == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	oIDInt, err := strconv.Atoi(oID)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is not a number: "+err.Error())))
		return
	}

	pID := chi.URLParam(r, paymentIdParam)

	if pID == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("paymentId is empty")))
		return
	}

	pIDInt, err := strconv.Atoi(pID)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("paymentId is not a number: "+err.Error())))
		return
	}

	err = order2.RemovePayment(r.Context(), uint(oIDInt), uint(pIDInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}

func (u *Router) GetProduct(w http.ResponseWriter, r *http.Request) {

	oID := chi.URLParam(r, orderIdParam)

	if oID == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	oIDInt, err := strconv.Atoi(oID)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is not a number: "+err.Error())))
		return
	}

	pID := chi.URLParam(r, productIdParam)

	if pID == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	pIDInt, err := strconv.Atoi(pID)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is not a number: "+err.Error())))
		return
	}

	is, err := order2.GetItem(r.Context(), uint(oIDInt), uint(pIDInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ListItemResponse{convertOrderItems(is)})

}

func (u *Router) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	oID := chi.URLParam(r, orderIdParam)

	if oID == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is empty")))
		return
	}

	oIDInt, err := strconv.Atoi(oID)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("orderId is not a number: "+err.Error())))
		return
	}

	var request DeleteItemsRequest

	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	for _, item := range request.Items {
		item.OrderID = uint(oIDInt)
	}

	err = order2.RemoveItems(r.Context(), uint(oIDInt), request.Items)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}
