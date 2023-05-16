package stock

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd/stock"
	"hades_backend/app/cmd/store"
	"hades_backend/app/database"
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
	err := db.AutoMigrate(&stock.Stock{}, &stock.Item{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &Router{db: db}
}

func (u *Router) URL() string {
	return "/stock"
}

const storeIdParam = "store_id"
const stockIdParam = "stock_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Delete("/{stock_id}", u.Delete)
		r.Put("/{stock_id}", u.Update)

		r.Post("/store/{store_id}", u.Create)
		r.Get("/store/{store_id}", u.GetByStoreId)
	}
}

func (u *Router) Update(w http.ResponseWriter, r *http.Request) {

	stockId := chi.URLParam(r, stockIdParam)

	if stockId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("stockId is empty")))
		return
	}

	stockIdInt, err := strconv.Atoi(stockId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("stockId is not a number: "+err.Error())))
		return
	}

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	stockDb, err := stock.UpdateStock(r.Context(), uint(stockIdInt), data.Stock)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	s := entityToResponse(stockDb)
	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{s})
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	db := database.DB.WithContext(r.Context())
	storeDb, err := stock.CreateStock(r.Context(), db, data.Stock)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	//db save
	s := entityToResponse(storeDb)
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{s})
}

func (u *Router) GetByStoreId(w http.ResponseWriter, r *http.Request) {

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	storeIdInt, err := strconv.Atoi(storeId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is not a number: "+err.Error())))
		return
	}

	sDb, err := stock.GetStock(r.Context(), uint(storeIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	s := entityToResponse(sDb)

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{s})

}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

	stockId := chi.URLParam(r, stockIdParam)

	if stockId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("stockId is empty")))
		return
	}

	stockIdInt, err := strconv.Atoi(stockId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("stockId is not a number: "+err.Error())))
		return
	}

	err = stock.DeleteStock(r.Context(), uint(stockIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func entityToResponse(s *stock.Stock) *model.Stock {

	st := storeEntityToModel(s.StoreID, s.Store)

	items := itemsEntityToModel(s.Items)

	return &model.Stock{
		ID:           s.ID,
		Store:        st,
		LastModified: s.UpdatedAt.Format(time.RFC3339),
		Items:        items,
	}
}

func storeEntityToModel(storeId uint, s *store.Store) *model.Store {
	if s == nil {
		return &model.Store{
			ID: storeId,
		}
	}
	return &model.Store{
		ID:      s.ID,
		Name:    s.Name,
		Address: s.Address,
	}
}

func itemsEntityToModel(i []*stock.Item) []*model.StockItem {
	items := make([]*model.StockItem, 0, len(i))

	for _, item := range i {
		items = append(items, itemEntityToModel(item))
	}

	return items
}

func itemEntityToModel(i *stock.Item) *model.StockItem {

	return &model.StockItem{
		ProductID:   i.ProductID,
		ProductName: i.Product.Name,
		ImageUrl:    i.Product.ImageUrl,
		Current:     i.Current,
		Suggested:   i.Suggested,
		AvgPrice:    i.AvgPrice,
	}
}
