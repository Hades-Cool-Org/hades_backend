package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"hades_backend/app/cmd/product"
	"hades_backend/app/cmd/store"
	"hades_backend/app/cmd/user"
	"hades_backend/app/cmd/vendors"
	"hades_backend/app/database"
	"hades_backend/app/logging"
	"hades_backend/app/model/order"
	"net/url"
	"strconv"
	"sync"
)

type Status string

const (
	Created           Status = "CRIADO"
	Accepted          Status = "ACEITO"
	AcceptedPartially Status = "ACEITO_PARCIAL"
	Received          Status = "RECEBIDO"
	ReceivedPartially Status = "RECEBIDO_PARCIAL"
	Completed         Status = "COMPLETADO"
)

type Order struct {
	gorm.Model
	Vendor   *vendors.Vendor
	State    string
	User     *user.User
	Total    decimal.Decimal `gorm:"type:decimal(7,6);"`
	Payments []*Payment
	Items    []*Item

	ModificationLock sync.Mutex `json:"-" sql:"-" gorm:"-"`
}

func (o Order) TableName() string {
	return "orders"
}

type OrderPrices struct {
	Total          decimal.Decimal
	PendingPayment decimal.Decimal
	Payed          decimal.Decimal
}

func (o *Order) CalculatedTotal() *OrderPrices {

	var total = decimal.Zero
	var payed = decimal.Zero
	var pendingPayment = decimal.Zero

	for _, item := range o.Items {
		total = item.CalculateTotal().Add(total)
	}

	for _, payment := range o.Payments {
		payed = payment.Total.Add(payed)
	}

	pendingPayment = total.Sub(payed)

	return &OrderPrices{
		Total:          total,
		PendingPayment: pendingPayment,
		Payed:          payed,
	}
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.State = string(Created)
	o.Total = o.CalculatedTotal().Total
	return nil
}

func (o *Order) UpdateTotal() {
	o.Total = o.CalculatedTotal().Total
}

// UpInsertItems updates the items of the order
func (o *Order) UpInsertItems(newItems []*order.Item) {

	o.ModificationLock.Lock()
	defer o.ModificationLock.Unlock()

	mapItem := make(map[string]*Item)

	fnKey := func(storeId, productId uint) string {
		return strconv.Itoa(int(storeId)) + "#" + strconv.Itoa(int(productId))
	}

	for _, item := range o.Items {
		mapItem[fnKey(item.StoreID, item.ProductID)] = item
	}

	for _, item := range newItems {
		key := fnKey(item.StoreID, item.ProductID)
		if _, ok := mapItem[key]; ok {
			mapItem[key].Quantity = item.Quantity
			mapItem[key].Available = item.Available
		} else {
			o.Items = append(o.Items, &Item{
				OrderID:      o.ID,
				ProductID:    item.ProductID,
				StoreID:      item.StoreID,
				Quantity:     item.Quantity,
				Available:    item.Available,
				PricePerItem: item.CalculateUnitPrice(),
			})
		}
	}
}

// RemoveItems removes given items from the order
func RemoveItems(ctx context.Context, orderID uint, items []*order.Item) error {

	db := database.DB.WithContext(ctx)
	l := logging.FromContext(ctx)

	l.Info(fmt.Sprintf("Removing items from order [%d]", orderID))

	// verify that the orderParams exists
	existingOrder := new(Order)

	if err := orderQuery(db).First(existingOrder, "id = ?", orderID).Error; err != nil {
		return ParseMysqlError(ctx, "order", err)
	}

	existingOrder.ModificationLock.Lock()
	defer existingOrder.ModificationLock.Unlock()

	itemsToRemove := make([]*Item, 0)

	for _, item := range existingOrder.Items {
		i := &Item{OrderID: orderID, ProductID: item.ProductID, StoreID: item.StoreID}
		itemsToRemove = append(itemsToRemove, i)
	}

	if err := db.Delete(itemsToRemove).Error; err != nil {
		return ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Removed items from order [%d]", orderID))

	return nil
}

// UpdateOrder updates an order
// items are not removed if they are not in the new list
func UpdateOrder(ctx context.Context, orderID uint, orderParams order.Order) error {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	marshal, err := json.Marshal(orderParams)

	if err != nil {
		return err
	}

	l.Info(fmt.Sprintf("Updating orderParams [%d] -> \n [%s]", orderID, string(marshal)))

	// verify that the orderParams exists
	existingOrder := new(Order)

	if err := orderQuery(db).First(existingOrder, "id = ?", orderID).Error; err != nil {
		return ParseMysqlError(ctx, "order", err)
	}

	if orderParams.User.ID != 0 {
		u := new(user.User)
		if err := db.First(u, "id = ?", orderParams.User.ID).Error; err != nil {
			return ParseMysqlError(ctx, "order", err)
		}
		existingOrder.User = u
	}

	if orderParams.Vendor.ID != 0 {
		v := new(vendors.Vendor)
		if err := db.First(v, "id = ?", orderParams.Vendor.ID).Error; err != nil {
			return ParseMysqlError(ctx, "order", err)
		}
		existingOrder.Vendor = v
	}

	tx := db.Begin()

	if len(orderParams.Payments) > 0 {
		newPayments := make([]*Payment, 0)
		for _, payment := range orderParams.Payments {
			p := &Payment{
				Type:   payment.Type,
				Total:  payment.Total,
				OderID: orderID,
				Text:   payment.Text,
			}
			newPayments = append(newPayments, p)
		}

		if err := tx.Model(existingOrder).Association("Payments").Replace(newPayments); err != nil {
			tx.Rollback()
			return ParseMysqlError(ctx, "order", err)
		}
	}

	if len(orderParams.Items) > 0 {
		existingOrder.UpInsertItems(orderParams.Items)
		//does it remove the old items? (not in the new list)
		//if err := tx.Model(existingOrder).Association("Items").Replace(existingOrder.Items); err != nil {
		//	tx.Rollback()
		//	return ParseMysqlError(ctx, "order", err)
		//}
		prices := existingOrder.CalculatedTotal()
		existingOrder.Total = prices.Total
	}

	if err := tx.Save(existingOrder).Error; err != nil {
		tx.Rollback()
		return ParseMysqlError(ctx, "order", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Updated order [%d] ", orderID))

	return nil
}

// CreateOrder creates a new order
// no payments nor items are added
func CreateOrder(ctx context.Context, orderParams order.Order) (*Order, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	marshal, err := json.Marshal(orderParams)

	if err != nil {
		return nil, err
	}

	l.Info(fmt.Sprintf("Creating order -> \n [%s]", string(marshal)))

	o := &Order{}

	if orderParams.User.ID != 0 {
		u := new(user.User)
		if err := db.First(u, "id = ?", orderParams.User.ID).Error; err != nil {
			return nil, ParseMysqlError(ctx, "user", err)
		}
		o.User = u
	}

	if orderParams.Vendor.ID != 0 {
		v := new(vendors.Vendor)
		if err := db.First(v, "id = ?", orderParams.Vendor.ID).Error; err != nil {
			return nil, ParseMysqlError(ctx, "vendor", err)
		}
		o.Vendor = v
	}

	if err := db.Omit("Items").Omit("Payments").Create(o).Error; err != nil {
		return nil, ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Created order [%d] ", o.ID))

	return o, nil
}

// GetOrder returns an order
func GetOrder(ctx context.Context, orderID uint) (*Order, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	o := new(Order)

	if err := orderQuery(db).First(o, "id = ?", orderID).Error; err != nil {
		return nil, ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Found order [%d] ", o.ID))

	return o, nil
}

type GetOrdersOptions struct {
	Params url.Values

	Offset int
	Limit  int
}

// GetOrders returns all orders
func GetOrders(ctx context.Context, options *GetOrdersOptions) ([]*Order, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	orders := make([]*Order, 0)

	query := parseOrderParams(db, options.Params)

	// not fetching relations for now
	if err := query.Find(&orders).Error; err != nil {
		return nil, ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Found %d orders", len(orders)))

	return orders, nil
}

func parseOrderParams(query *gorm.DB, params url.Values) *gorm.DB {

	tableName := (&Order{}).TableName()

	if s := params.Get("status"); s != "" {
		query = query.Where(tableName+".status = ?", s)
	}

	if s := params.Get("vendor_id"); s != "" {
		query = query.Where(tableName+".vendor.id = ?", s)
	}

	if s := params.Get("user_id"); s != "" {
		query = query.Where(tableName+".user.id = ?", s)
	}

	if s := params.Get("created_at"); s != "" {
		//TODO: parse date
		query = query.Where(tableName+".created_at >= ?", s)
	}

	return query
}

type Payment struct {
	gorm.Model
	Type   string
	Total  decimal.Decimal `gorm:"type:decimal(7,6);"`
	OderID uint
	Text   string `gorm:"type:text"`
}

func (p Payment) TableName() string {
	return "payments"
}

type Item struct {
	OrderID uint `gorm:"primaryKey;autoIncrement:false"`

	ProductID uint `gorm:"primaryKey;autoIncrement:false"`
	Product   *product.Product

	StoreID uint `gorm:"primaryKey;autoIncrement:false"`
	Store   *store.Store

	Quantity     float64
	Available    float64
	PricePerItem decimal.Decimal `gorm:"type:decimal(7,6);"`
}

func (i *Item) CalculateTotal() decimal.Decimal {
	return i.PricePerItem.Mul(decimal.NewFromFloat(i.Quantity))
}

func (i Item) TableName() string {
	return "items"
}

func orderQuery(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Payments").
		Preload("Items").
		Preload("Vendor").
		Preload("User")
}
