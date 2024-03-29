package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd"
	"hades_backend/app/cmd/product"
	"hades_backend/app/cmd/store"
	"hades_backend/app/database"
	"hades_backend/app/logging"
	"hades_backend/app/model"
	"net/url"
)

type Stock struct {
	gorm.Model

	StoreID uint `gorm:"index;not null;"`
	Store   *store.Store

	Items []*Item
}

// TableName overrides the table name used by StockItem to `stock_products`
func (Stock) TableName() string {
	return "stock"
}

type Item struct {
	StockID uint `gorm:"primaryKey;autoIncrement:false"`

	ProductID uint `gorm:"primaryKey;autoIncrement:false"`
	Product   *product.Product

	Current   decimal.Decimal
	Suggested decimal.Decimal

	AvgPrice decimal.Decimal `gorm:"type:decimal(12,3);"`
}

func (s *Stock) BeforeDelete(tx *gorm.DB) error {
	delModels := map[string]interface{}{
		"items": &[]Item{},
	}
	for name, dm := range delModels {
		if result := tx.Delete(dm, "stock_id = ?", s.ID); result.Error != nil {
			return errors.Wrap(result.Error, fmt.Sprintf("Error deleting %s records", name))
		}
	}
	return nil
}

// TableName overrides the table name used by StockItem to `stock_products`
func (Item) TableName() string {
	return "stock_items"
}

func (s *Stock) findItem(productID uint) *Item {
	for _, item := range s.Items {
		if item.ProductID == productID {
			return item
		}
	}
	return nil
}

// GetStockOptions TODO: add pagination
type GetStockOptions struct {
	Params url.Values
}

func (o *GetStockOptions) parseStockParams(query *gorm.DB) *gorm.DB {

	tableName := (&Stock{}).TableName()

	if s := o.Params.Get("id"); s != "" {
		query = query.Where(tableName+".id = ?", s)
	}

	if s := o.Params.Get("store_id"); s != "" {
		query = query.Where(tableName+".store_id = ?", s)
	}

	return query
}

func GetStock(ctx context.Context, storeID uint) (*Stock, error) {
	db := database.DB.WithContext(ctx)
	l := logging.FromContext(ctx)

	l.Info(fmt.Sprintf("Getting stock for store %d", storeID))

	s := &Stock{}

	if result := db.
		Preload("Items").
		Preload("Items.Product").
		First(s, "store_id = ?", storeID); result.Error != nil {
		return nil, cmd.ParseMysqlError(ctx, "stock", result.Error)
	}

	//if result := db.Model(s).Association("Items").Find(&s.Items); result.Error != nil {
	//	return nil, result.Error
	//}

	return s, nil
}

func GetStocks(ctx context.Context, opts *GetStockOptions) ([]*Stock, error) {
	db := database.DB.WithContext(ctx)
	l := logging.FromContext(ctx)

	l.Info(fmt.Sprintf("Getting stock for store %v", opts))

	query := db.Preload("Items").Preload("Items.Product")

	query = opts.parseStockParams(query)

	stocks := make([]*Stock, 0)

	if err := query.Find(&stocks).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "occurrence", err)
	}

	return stocks, nil
}

func CreateStock(ctx context.Context, db *gorm.DB, stockParams *model.Stock) (*Stock, error) {

	l := logging.FromContext(ctx)

	marshal, err := json.Marshal(stockParams)

	if err != nil {
		return nil, err
	}

	l.Info(fmt.Sprintf("Creating stock -> \n [%s]", string(marshal)))

	s := &Stock{}

	if stockParams.Store.ID != 0 {
		st := new(store.Store)
		if result := db.First(st, stockParams.Store.ID); result.Error != nil {
			return nil, result.Error
		}
		s.Store = st
		s.StoreID = st.ID
	} else {
		return nil, net.NewBadRequestError(ctx, errors.New("storeID is required"))
	}

	items := make([]*Item, 0)
	for _, item := range stockParams.Items {
		if item.ProductID != 0 {
			p := new(product.Product)
			if result := db.First(p, item.ProductID); result.Error != nil {
				return nil, result.Error
			}
			items = append(items, &Item{
				Product:   p,
				ProductID: p.ID,
				Current:   item.Current,
				Suggested: item.Suggested,
				AvgPrice:  item.AvgPrice,
			})
		} else {
			l.Info(fmt.Sprintf("missing product id... skipping [%v]", item))
		}
	}

	s.Items = items

	if result := db.Create(s); result.Error != nil {
		return nil, cmd.ParseMysqlError(ctx, "stock", result.Error)
	}

	return s, nil
}

func UpdateStock(ctx context.Context, stockID uint, stockParams *model.Stock) (*Stock, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	marshal, err := json.Marshal(stockParams)

	if err != nil {
		return nil, err
	}

	l.Info(fmt.Sprintf("Updating stock -> \n [%s]", string(marshal)))

	s := &Stock{}

	if result := db.Preload("Items").First(s, stockID); result.Error != nil {
		return nil, cmd.ParseMysqlError(ctx, "stock", result.Error)
	}

	if len(stockParams.Items) > 0 {
		err := s.doAction(ctx, ActionTypeUpsert, stockParams.Items)
		if err != nil {
			return nil, cmd.ParseMysqlError(ctx, "stock", err)
		}
	}

	if result := db.Save(s); result.Error != nil {
		return nil, cmd.ParseMysqlError(ctx, "stock", result.Error)
	}

	return s, nil
}

func AddStockItem(ctx context.Context, db *gorm.DB, stockID uint, itemParams []*model.StockItem) error {

	l := logging.FromContext(ctx)

	marshal, err := json.Marshal(itemParams)

	if err != nil {
		return err
	}

	l.Info(fmt.Sprintf("Adding stock item -> \n [%s]", string(marshal)))

	s := &Stock{}

	if result := db.Preload("Items").First(s, stockID); result.Error != nil {
		return cmd.ParseMysqlError(ctx, "stock", result.Error)
	}

	if len(itemParams) > 0 {
		err := s.doAction(ctx, ActionTypeAddition, itemParams)
		if err != nil {
			return cmd.ParseMysqlError(ctx, "stock", err)
		}
	}

	if result := db.Save(s); result.Error != nil {
		return cmd.ParseMysqlError(ctx, "stock", result.Error)
	}

	return nil
}

func SubtractStockItem(ctx context.Context, stockID uint, itemParams []*model.StockItem) error {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	marshal, err := json.Marshal(itemParams)

	if err != nil {
		return err
	}

	l.Info(fmt.Sprintf("Subtracting stock item -> \n [%s]", string(marshal)))

	s := &Stock{}

	if result := db.Preload("Items").First(s, stockID); result.Error != nil {
		return cmd.ParseMysqlError(ctx, "stock", result.Error)
	}

	if len(itemParams) > 0 {
		err := s.doAction(ctx, ActionTypeSubtraction, itemParams)
		if err != nil {
			return cmd.ParseMysqlError(ctx, "stock", err)
		}
	}

	if result := db.Save(s); result.Error != nil {
		return cmd.ParseMysqlError(ctx, "stock", result.Error)
	}

	return nil
}

func DeleteStock(ctx context.Context, stockID uint) error {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Deleting stock -> \n [%d]", stockID))

	s := &Stock{}

	if result := db.First(s, stockID); result.Error != nil {
		return cmd.ParseMysqlError(ctx, "stock", result.Error)
	}

	if result := db.Delete(s); result.Error != nil {
		return cmd.ParseMysqlError(ctx, "stock", result.Error)
	}

	return nil
}

type ActionType int

const (
	ActionTypeUpsert      ActionType = iota // atualiza tudo que for informado e adiciona o que não existir
	ActionTypeAddition                      // adiciona a quantidade informado ao estoque atual
	ActionTypeSubtraction                   // subtrai a quantidade informado ao estoque atual
)

func (s *Stock) upsertStock(ctx context.Context, newItems []*model.StockItem) error {
	fnKey := func(stockId, productId uint) string {
		return fmt.Sprintf("%d#%d", stockId, productId)
	}

	oldItems := make(map[string]*Item)

	for _, item := range s.Items {
		oldItems[fnKey(s.ID, item.ProductID)] = item
	}

	i := make([]*Item, 0)

	for _, item := range newItems {
		key := fnKey(s.ID, item.ProductID)

		if oldItem, ok := oldItems[key]; ok {
			oldItem.Current = item.Current
			oldItem.Suggested = item.Suggested
			oldItem.AvgPrice = item.AvgPrice
			i = append(i, oldItem)
		} else {
			i = append(i, &Item{
				StockID:   s.ID,
				ProductID: item.ProductID,
				Current:   item.Current,
				Suggested: item.Suggested,
				AvgPrice:  item.AvgPrice,
			})
		}
	}

	s.Items = i
	return nil
}

func calculateAvgPrice(currentItem *Item, newItem *model.StockItem) decimal.Decimal {

	qty := currentItem.Current.Add(newItem.Current)

	total := currentItem.AvgPrice.
		Mul(currentItem.Current).
		Add(newItem.AvgPrice.Mul(newItem.Current))

	return total.Div(qty)
}

func (s *Stock) addStock(ctx context.Context, items []*model.StockItem) error {
	for _, item := range items {
		if i := s.findItem(item.ProductID); i != nil {
			i.Current = i.Current.Add(item.Current)
			i.AvgPrice = calculateAvgPrice(i, item)
		} else {
			s.Items = append(s.Items, &Item{
				StockID:   s.ID,
				ProductID: item.ProductID,
				Current:   item.Current,
				Suggested: item.Suggested,
				AvgPrice:  item.AvgPrice,
			})
		}
	}
	return nil
}

func (s *Stock) subtractStock(ctx context.Context, items []*model.StockItem) error {
	for _, item := range items {
		if i := s.findItem(item.ProductID); i != nil {
			i.Current = i.Current.Sub(item.Current)
		}
	}
	return nil
}

func (s *Stock) doAction(ctx context.Context, t ActionType, items []*model.StockItem) error {
	actions := map[ActionType]func(ctx context.Context, items []*model.StockItem) error{
		ActionTypeUpsert:      s.upsertStock,
		ActionTypeAddition:    s.addStock,
		ActionTypeSubtraction: s.subtractStock,
	}

	err := actions[t](ctx, items)
	if err != nil {
		return err
	}

	return nil
}
