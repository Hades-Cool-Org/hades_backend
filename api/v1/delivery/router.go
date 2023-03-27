package delivery

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/api/utils/net"
	"net/http"
	"time"
)

type Router struct {
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

		r.Post("/turn", u.CreateUserTurn)
		r.Get("/turn", u.GetAllTurns)
		r.Post("/turn/{turn_id}/end", u.EndUserTurn)
		r.Delete("/turn/{turn_id}", u.DeleteUserTurn)
		r.Get("/turn/{turn_id}", u.GetUserTurn)

		r.Post("/vehicles", u.CreateVehicle)                 //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
		r.Get("/vehicles", u.GetAllVehicles)                 //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
		r.Get("/vehicles/{vehicle_id}", u.GetVehicle)        //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
		r.Delete("/vehicles/{vehicle_id}", u.DeleteVehicles) //Associar um carro a um entregador //todo: usar mesma funcao que end delivery")
	}
}
