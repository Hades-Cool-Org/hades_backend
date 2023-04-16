package delivery

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd/delivery"
	"hades_backend/app/model"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Router struct {
	db *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	err := db.AutoMigrate(&delivery.Delivery{}, &delivery.Item{}, &delivery.Vehicle{}, &delivery.Session{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &Router{db: db}
}

func (u *Router) URL() string {
	return "/deliveries"
}

const deliveryIdParam = "delivery_id"
const userIdParam = "user_id"
const sessionIdParam = "session_id"
const vehicleIdParam = "vehicle_id"

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
		r.Post("/sessions", u.CreateSession)
		r.Get("/sessions", u.GetAllSessions)
		r.Delete("/sessions/{session_id}", u.EndSession)
		r.Get("/sessions/{session_id}", u.GetSession)
		//
		r.Post("/vehicles", u.CreateVehicle)                //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
		r.Get("/vehicles", u.GetAllVehicles)                //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
		r.Get("/vehicles/{vehicle_id}", u.GetVehicle)       //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
		r.Delete("/vehicles/{vehicle_id}", u.DeleteVehicle) //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
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

func (u *Router) CreateSession(w http.ResponseWriter, r *http.Request) {

	var request SessionRequest

	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	s, err := delivery.CreateSession(r.Context(), request.Session)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &SessionResponse{convertSessionToResponse(s)})

}

func (u *Router) GetAllSessions(w http.ResponseWriter, r *http.Request) {

	var opts *delivery.GetSessionOptions

	opts.Params = r.URL.Query()

	sessions, err := delivery.GetSessions(r.Context(), opts)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &SessionListResponse{convertMultipleSessionsToResponse(sessions)})
}

func (u *Router) GetSession(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, sessionIdParam)

	if param == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("sessionId is empty")))
		return
	}

	sessionIdInt, err := strconv.Atoi(param)

	s, err := delivery.GetSession(r.Context(), uint(sessionIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &SessionResponse{convertSessionToResponse(s)})
}

func (u *Router) EndSession(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, sessionIdParam)

	if param == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("sessionId is empty")))
		return
	}

	sessionIdInt, err := strconv.Atoi(param)

	err = delivery.DeleteSession(r.Context(), uint(sessionIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (u *Router) CreateVehicle(w http.ResponseWriter, r *http.Request) {

	var request VehicleRequest

	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	v, err := delivery.CreateVehicle(r.Context(), request.Vehicle)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &VehicleResponse{convertVehicleToResponse(v)})
}

func (u *Router) GetAllVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := delivery.GetAllVehicles(r.Context())

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ListVehicleResponse{convertMultipleVehiclesToResponse(vehicles)})
}

func (u *Router) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, vehicleIdParam)

	if param == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("vehicleId is empty")))
		return
	}

	vehicleIdInt, err := strconv.Atoi(param)

	err = delivery.DeleteVehicle(r.Context(), uint(vehicleIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (u *Router) GetVehicle(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, vehicleIdParam)

	if param == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("vehicleId is empty")))
		return
	}

	vehicleIdInt, err := strconv.Atoi(param)

	v, err := delivery.GetVehicle(r.Context(), uint(vehicleIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &VehicleResponse{convertVehicleToResponse(v)})
}

func convertMultipleVehiclesToResponse(vehicles []*delivery.Vehicle) []*model.Vehicle {
	var responses []*model.Vehicle
	for _, v := range vehicles {
		responses = append(responses, convertVehicleToResponse(v))
	}
	return responses
}

func convertVehicleToResponse(v *delivery.Vehicle) *model.Vehicle {
	return &model.Vehicle{
		ID:   v.ID,
		Name: v.Name,
		Type: v.Type,
	}
}

func convertMultipleDeliveriesToResponse(deliveries []*delivery.Delivery) []*model.Delivery {
	var responses []*model.Delivery
	for _, d := range deliveries {
		responses = append(responses, convertDeliveryToResponse(d))
	}
	return responses
}

func convertMultipleSessionsToResponse(sessions []*delivery.Session) []*model.Session {
	var responses []*model.Session
	for _, s := range sessions {
		responses = append(responses, convertSessionToResponse(s))
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
