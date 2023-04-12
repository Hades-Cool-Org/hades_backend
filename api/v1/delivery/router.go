package delivery

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd/delivery"
	"hades_backend/app/model"
	"net/http"
	"strconv"
	"time"
)

type Router struct {
	db *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	err := db.AutoMigrate(&delivery.Delivery{}, &delivery.Item{}, &delivery.Vehicle{}, &delivery.Session{})
	if err != nil {
		panic(err)
	}
	return &Router{db: db}
}

func (u *Router) URL() string {
	return "/delivery"
}

const deliveryIdParam = "delivery_id"
const userIdParam = "user_id"
const dateStartParam = "date_start"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Create)
		r.Get("/", u.GetAll)
		r.Delete("/{delivery_id}", u.Delete)
		r.Get("/{delivery_id}", u.Get)
		r.Put("/{delivery_id}", u.Update)
		r.Delete("/{delivery_id}/items", u.RemoveItems)

		//// gerente na loja conferindo pedido// TODO MELHOR NOME?
		//r.Post("/{delivery_id}/conference", u.KKKK)
		//
		//r.Post("/session", u.CreateUserTurn)
		//r.Get("/session", u.GetAllTurns)
		//r.Post("/session/{session_id}/end", u.EndUserTurn)
		//r.Delete("/session/{session_id}", u.DeleteUserTurn)
		//r.Get("/session/{session_id}", u.GetUserTurn)
		//
		//r.Post("/vehicles", u.CreateVehicle)                 //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
		//r.Get("/vehicles", u.GetAllVehicles)                 //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
		//r.Get("/vehicles/{vehicle_id}", u.GetVehicle)        //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
		//r.Delete("/vehicles/{vehicle_id}", u.DeleteVehicles) //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
	}
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {
	var request Request

	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	d, err := delivery.CreateDelivery(r.Context(), request.Delivery)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{convertDeliveryToResponse(d)})
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {
	opts := &delivery.GetDeliveryOptions{
		Params: r.URL.Query(),
	}

	deliveries, err := delivery.GetDeliveries(r.Context(), opts)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ListResponse{convertMultipleDeliveriesToResponse(deliveries)})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {
	deliveryId := chi.URLParam(r, deliveryIdParam)

	if deliveryId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is empty")))
		return
	}

	deliveryIdInt, err := strconv.Atoi(deliveryId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is not a number: "+err.Error())))
		return
	}

	err = delivery.DeleteDelivery(r.Context(), uint(deliveryIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.Render(w, r, net.NoContent())
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {
	deliveryId := chi.URLParam(r, deliveryIdParam)

	if deliveryId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is empty")))
		return
	}

	deliveryIdInt, err := strconv.Atoi(deliveryId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is not a number: "+err.Error())))
		return
	}

	d, err := delivery.GetDelivery(r.Context(), uint(deliveryIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{convertDeliveryToResponse(d)})
}

func (u *Router) Update(w http.ResponseWriter, r *http.Request) {
	deliveryId := chi.URLParam(r, deliveryIdParam)

	if deliveryId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is empty")))
		return
	}

	deliveryIdInt, err := strconv.Atoi(deliveryId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is not a number: "+err.Error())))
		return
	}

	var request Request

	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	d, err := delivery.UpdateDelivery(r.Context(), uint(deliveryIdInt), request.Delivery)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{convertDeliveryToResponse(d)})
}

func (u *Router) RemoveItems(w http.ResponseWriter, r *http.Request) {
	deliveryId := chi.URLParam(r, deliveryIdParam)

	if deliveryId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is empty")))
		return
	}

	deliveryIdInt, err := strconv.Atoi(deliveryId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("deliveryId is not a number: "+err.Error())))
		return
	}

	var request ItemRequest

	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	err = delivery.RemoveDeliveryItems(r.Context(), uint(deliveryIdInt), request.DeliveryItems)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}

func convertMultipleDeliveriesToResponse(deliveries []*delivery.Delivery) []*model.Delivery {
	var responses []*model.Delivery
	for _, d := range deliveries {
		responses = append(responses, convertDeliveryToResponse(d))
	}
	return responses
}

func convertSessionToResponse(s *delivery.Session) *model.Session {
	return &model.Session{
		ID: s.ID,
		User: &model.User{
			ID:    s.UserID,
			Name:  s.User.Name,
			Email: s.User.Email,
			Phone: s.User.Phone,
		},
		Vehicle: &model.Vehicle{
			ID:   s.VehicleID,
			Name: s.Vehicle.Name,
			Type: s.Vehicle.Type,
		},
		StartDate: s.CreatedAt.Format(time.RFC3339),
		EndDate:   s.EndDate.Time.Format(time.RFC3339),
	}
}

func convertItemsToResponse(items []*delivery.Item) []*model.DeliveryItem {

	var deliveryItems []*model.DeliveryItem

	for _, i := range items {
		deliveryItems = append(deliveryItems, &model.DeliveryItem{
			ProductID:     i.ProductID,
			StoreID:       i.StoreID,
			Name:          i.Product.Name,
			ImageUrl:      i.Product.ImageUrl,
			MeasuringUnit: i.Product.MeasuringUnit,
			Quantity:      i.Quantity,
		})
	}
	return deliveryItems
}

func convertDeliveryToResponse(d *delivery.Delivery) *model.Delivery {

	deliveryState, _ := model.DeliveryStateFromString(d.State)

	orderState, _ := model.OrderStateFromString(d.Order.State)

	o := &model.Order{
		ID:          d.OrderID,
		Vendor:      nil,
		CreatedDate: d.Order.CreatedAt.Format(time.RFC3339),
		State:       orderState,
		EndDate:     d.Order.CompletedDate.Time.Format(time.RFC3339),
		User: &model.User{
			ID:    d.Order.UserID,
			Name:  d.Order.User.Name,
			Email: d.Order.User.Email,
			Phone: d.Order.User.Phone,
		},
		Total: decimal.Decimal{},
		Payments: func() []*model.Payment {
			var payments []*model.Payment
			for _, p := range d.Order.Payments {
				payments = append(payments, &model.Payment{
					ID:    p.ID,
					Type:  p.Type,
					Total: p.Total,
					Date:  p.CreatedAt.Format(time.RFC3339),
					Text:  p.Text,
				})
			}
			return payments
		}(),
		Items: nil, //items empty, not sure if we will need that
	}

	s := convertSessionToResponse(d.Session)

	m := &model.Delivery{
		ID:            d.ID,
		State:         &deliveryState,
		StartDate:     d.CreatedAt.Format(time.RFC3339),
		EndDate:       d.EndDate.Time.Format(time.RFC3339),
		Order:         o,
		Session:       s,
		DeliveryItems: convertItemsToResponse(d.Items),
	}

	return m
}
