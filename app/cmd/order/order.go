package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd"
	"hades_backend/app/cmd/balance"
	"hades_backend/app/cmd/user"
	"hades_backend/app/cmd/vendors"
	"hades_backend/app/database"
	"hades_backend/app/logging"
	"hades_backend/app/model"
	"net/url"
	"strconv"
	"sync"
)

type Order struct {
	gorm.Model

	VendorID uint `gorm:"index"`
	Vendor   *vendors.Vendor

	UserID uint `gorm:"index"`
	User   *user.User

	State string
	Total decimal.Decimal `gorm:"type:decimal(12,3);"`

	Payments []*Payment
	Items    []*Item

	CompletedDate sql.NullTime //TODO: mudança de estado setar isso

	ModificationLock sync.Mutex `json:"-" sql:"-" gorm:"-"`
}

func (o *Order) TableName() string {
	return "orders"
}

type Prices struct {
	Total          decimal.Decimal
	PendingPayment decimal.Decimal
	Payed          decimal.Decimal
}

func (o *Order) CalculatedTotal() *Prices {

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

	return &Prices{
		Total:          total.Round(2),
		PendingPayment: pendingPayment,
		Payed:          payed,
	}
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.State = string(model.Created)
	o.Total = o.CalculatedTotal().Total
	return nil
}

func (o *Order) BeforeDelete(tx *gorm.DB) error {

	delModels := map[string]interface{}{
		"items":    &[]Item{},
		"payments": &[]Payment{},
		//"delivery":   &[]delivery.Delivery{},
		//"occurrence": &[]occurence.Occurrence{},
	}
	for name, dm := range delModels {
		if result := tx.Delete(dm, "order_id = ?", o.ID); result.Error != nil {
			return errors.Wrap(result.Error, fmt.Sprintf("Error deleting %s records", name))
		}
	}
	return nil
}

func (o *Order) UpdateTotal() {
	o.Total = o.CalculatedTotal().Total
}

// updateItems updates the items of the order - NO ITEMS ARE REMOVED
func (o *Order) updateItems(newItems []*model.OrderItem) {

	o.ModificationLock.Lock()
	defer o.ModificationLock.Unlock()

	mapItem := make(map[string]*Item)

	fnKey := func(storeId, productId uint) string {
		return strconv.Itoa(int(storeId)) + "#" + strconv.Itoa(int(productId))
	}

	for _, item := range o.Items {
		mapItem[fnKey(item.StoreID, item.ProductID)] = item
	}

	i := make([]*Item, 0)

	for _, item := range newItems {
		key := fnKey(item.StoreID, item.ProductID)

		item.Quantity = item.Quantity.Round(3)
		item.Total = item.Total.Round(2)

		p := item.CalculateUnitPrice()

		if _, ok := mapItem[key]; ok {
			mapItem[key].Quantity = item.Quantity
			mapItem[key].UnitPrice = p
		} else {
			i = append(i, &Item{
				OrderID:   o.ID,
				ProductID: item.ProductID,
				StoreID:   item.StoreID,
				Quantity:  item.Quantity,
				UnitPrice: p,
			})
		}
	}

	for _, item := range mapItem {
		i = append(i, item)
	}

	o.Items = i
}

// RemoveItems removes given items from the order
func RemoveItems(ctx context.Context, orderID uint, items []*model.OrderItem) error {

	db := database.DB.WithContext(ctx)
	l := logging.FromContext(ctx)

	l.Info(fmt.Sprintf("Removing items from order [%d]", orderID))

	// verify that the orderParams exists
	existingOrder := new(Order)

	if err := orderQuery(db).First(existingOrder, "id = ?", orderID).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "order", err)
	}

	existingOrder.ModificationLock.Lock()
	defer existingOrder.ModificationLock.Unlock()

	itemsToRemove := make([]*Item, 0)

	for _, item := range items {
		i := &Item{OrderID: orderID, ProductID: item.ProductID, StoreID: item.StoreID}
		itemsToRemove = append(itemsToRemove, i)
	}

	if err := db.Delete(itemsToRemove).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Removed items from order [%d]", orderID))

	return nil
}

// UpdateOrder updates an order
// items are not removed if they are not in the new list
func UpdateOrder(ctx context.Context, orderID uint, orderParams *model.Order) error {

	db := database.DB.WithContext(ctx)

	return UpdateOrderInTx(ctx, db, orderID, orderParams)
}

// UpdateOrderInTx updates an order in transaction
// items are not removed if they are not in the new list
func UpdateOrderInTx(ctx context.Context, db *gorm.DB, orderID uint, orderParams *model.Order) error {

	l := logging.FromContext(ctx)

	marshal, err := json.Marshal(orderParams)

	if err != nil {
		return err
	}

	l.Info(fmt.Sprintf("Updating orderParams [%d] -> \n [%s]", orderID, string(marshal)))

	// verify that the orderParams exists
	existingOrder := new(Order)

	if err := orderQuery(db).First(existingOrder, "id = ?", orderID).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "order", err)
	}

	if orderParams.User != nil && orderParams.User.ID != 0 {
		u := new(user.User)
		if err := db.First(u, "id = ?", orderParams.User.ID).Error; err != nil {
			return cmd.ParseMysqlError(ctx, "user", err)
		}
		existingOrder.User = u
	}

	if orderParams.Vendor != nil && orderParams.Vendor.ID != 0 {
		v := new(vendors.Vendor)
		if err := db.First(v, "id = ?", orderParams.Vendor.ID).Error; err != nil {
			return cmd.ParseMysqlError(ctx, "vendor", err)
		}
		existingOrder.Vendor = v
	}

	tx := db.Begin()

	if len(orderParams.Payments) > 0 {
		newPayments := make([]*Payment, 0)
		for _, payment := range orderParams.Payments {
			p := &Payment{
				Type:    payment.Type,
				Total:   payment.Total,
				OrderID: orderID,
				Text:    payment.Text,
			}
			newPayments = append(newPayments, p)
		}

		if err := tx.Model(existingOrder).Association("Payments").Replace(newPayments); err != nil {
			tx.Rollback()
			return cmd.ParseMysqlError(ctx, "order", err)
		}
	}

	if len(orderParams.Items) > 0 {
		existingOrder.updateItems(orderParams.Items)
		//does it remove the old items? (not in the new list)
		//if err := tx.Model(existingOrder).Association("Items").Sa(existingOrder.Items); err != nil {
		//	tx.Rollback()
		//	return cmd.ParseMysqlError(ctx, "order", err)
		//}
		prices := existingOrder.CalculatedTotal()
		existingOrder.Total = prices.Total

		if orderParams.User == nil || orderParams.User.ID == 0 {
			tx.Rollback()
			return net.NewBadRequestError(ctx, errors.New("userId is required when updating items"))
		}

		err := validateBalance(ctx, tx, orderParams.User.ID, existingOrder)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if orderParams.State != nil {
		existingOrder.State = string(*orderParams.State) //TODO: not sure if we need to always update this
	}

	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(existingOrder).Error; err != nil {
		tx.Rollback()
		return cmd.ParseMysqlError(ctx, "order", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return cmd.ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Updated order [%d] ", orderID))

	return nil
}

func validateBalance(ctx context.Context, tx *gorm.DB, userID uint, existingOrder *Order) error {
	//validating user credit
	b, err := balance.GetBalance(ctx, tx, userID)

	if err != nil {
		tx.Rollback()
		return cmd.ParseMysqlError(ctx, "order", err)
	}

	if b == nil || b.Amount.LessThan(existingOrder.Total) {
		tx.Rollback()
		return net.NewBadRequestError(ctx, errors.New("user has not enough balance"))
	}

	_, err = balance.ManageBalance(ctx, &balance.Params{
		UserID:    userID,
		Amount:    existingOrder.Total,
		Operation: balance.Debit,
	})

	if err != nil {
		tx.Rollback()
		return net.NewBadRequestError(ctx, err)
	}

	return nil
}

// CreateOrder creates a new order
// no payments nor items are added
func CreateOrder(ctx context.Context, orderParams *model.Order) (*Order, error) {

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
			return nil, cmd.ParseMysqlError(ctx, "user", err)
		}
		o.User = u
	}

	if orderParams.Vendor.ID != 0 {
		v := new(vendors.Vendor)
		if err := db.First(v, "id = ?", orderParams.Vendor.ID).Error; err != nil {
			return nil, cmd.ParseMysqlError(ctx, "vendor", err)
		}
		o.Vendor = v
	}

	o.State = string(model.Created)

	if err := db.Omit("Items").Omit("Payments").Create(o).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "order", err)
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
		return nil, cmd.ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Found order [%d] ", o.ID))

	return o, nil
}

func GetItem(ctx context.Context, orderID, itemID uint) ([]*Item, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	i := new([]*Item)

	if err := db.Find(i, "order_id = ? AND product_id = ?", orderID, itemID).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Found item [%d] in order [%d]", itemID, orderID))

	return *i, nil
}

// GetOrdersOptions TODO: add pagination
type GetOrdersOptions struct {
	Params url.Values
}

// GetOrders returns all orders TODO: add pagination
func GetOrders(ctx context.Context, options *GetOrdersOptions) ([]*Order, error) {

	l := logging.FromContext(ctx)
	q := database.DB.WithContext(ctx)

	orders := make([]*Order, 0)

	query := orderQuery(q)
	query = options.parseOrderParams(query)

	// not fetching relations for now
	if err := query.Find(&orders).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Found %d orders", len(orders)))

	return orders, nil
}

func AddPayment(ctx context.Context, orderID uint, paymentParams *model.Payment) (*Payment, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	o := new(Order)

	if err := db.First(o, "id = ?", orderID).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "order", err)
	}

	p := &Payment{
		Type:    paymentParams.Type,
		Total:   paymentParams.Total,
		OrderID: orderID,
		Text:    paymentParams.Text,
	}

	if err := db.Create(p).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "payment", err)
	}

	l.Info(fmt.Sprintf("Added payment [%d] to order [%d]", p.ID, orderID))

	return p, nil
}

func RemovePayment(ctx context.Context, orderID uint, paymentID uint) error {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	p := new(Payment)

	if err := db.First(p, "id = ?", paymentID).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "payment", err)
	}

	if err := db.Delete(p).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "payment", err)
	}

	l.Info(fmt.Sprintf("Removed payment [%d] from order [%d]", p.ID, orderID))

	return nil
}

func DeleteOrder(ctx context.Context, orderID uint) error {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	o := new(Order)

	if err := db.First(o, "id = ?", orderID).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "order", err)
	}

	if err := db.Delete(o).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "order", err)
	}

	l.Info(fmt.Sprintf("Deleted order [%d]", orderID))

	return nil
}

func (o *GetOrdersOptions) parseOrderParams(query *gorm.DB) *gorm.DB {

	tableName := (&Order{}).TableName()

	if s := o.Params.Get("state"); s != "" {
		query = query.Where(tableName+".state = ?", s)
	}

	if s := o.Params.Get("vendor_id"); s != "" {
		query = query.Where(tableName+".vendor_id = ?", s)
	}

	if s := o.Params.Get("user_id"); s != "" {
		query = query.Where(tableName+".user_id = ?", s)
	}

	if s := o.Params.Get("created_at"); s != "" {
		//TODO: parse date
		query = query.Where(tableName+".created_at >= ?", s)
	}

	return query
}

func orderQuery(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Payments").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Store").
		Preload("Vendor").
		Preload("User")
}
