package delivery

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
	"hades_backend/app/cmd/order"
	"hades_backend/app/cmd/product"
	"hades_backend/app/cmd/store"
	"hades_backend/app/database"
	"hades_backend/app/logging"
	"hades_backend/app/model"
	"net/http"
	"net/url"
	"time"
)

type Delivery struct {
	gorm.Model
	State string //ABERTO,COLETADO,ENTREGUE

	EndDate sql.NullTime `gorm:"index"`

	SessionID uint
	Session   *Session

	OrderID uint `gorm:"index"`
	Order   *order.Order

	Items []*Item
}

func (d Delivery) TableName() string {
	return "deliveries"
}

func (d *Delivery) BeforeDelete(tx *gorm.DB) error {

	delModels := map[string]interface{}{
		"items": &[]Item{},
	}
	for name, dm := range delModels {
		if result := tx.Delete(dm, "delivery_id = ?", d.ID); result.Error != nil {
			return errors.Wrap(result.Error, fmt.Sprintf("Error deleting %s records", name))
		}
	}
	return nil
}

func (d *Delivery) BeforeCreate(tx *gorm.DB) (err error) {
	d.State = string(model.OPENED)
	return nil
}

type Item struct {
	DeliveryID uint `gorm:"primaryKey;autoIncrement:false"`

	ProductID uint `gorm:"primaryKey;autoIncrement:false"`
	Product   *product.Product

	StoreID uint `gorm:"primaryKey;autoIncrement:false"`
	Store   *store.Store

	UnitPrice decimal.Decimal `gorm:"type:decimal(12,3);"`

	Quantity float64
}

// CreateDelivery creates a new delivery
func CreateDelivery(ctx context.Context, deliveryParam *model.Delivery) (*Delivery, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	marshal, err := json.Marshal(deliveryParam)

	if err != nil {
		return nil, err
	}

	l.Info(fmt.Sprintf("Creating delivery -> \n [%s]", string(marshal)))

	d := &Delivery{}

	if deliveryParam.Session != nil {
		if deliveryParam.Session.ID != 0 {
			s := new(Session)
			if err := db.
				Preload("Vehicle").
				Preload("User").
				First(s, "id = ?", deliveryParam.Session.ID).Error; err != nil {
				return nil, cmd.ParseMysqlError(ctx, "session", err)
			}
			d.Session = s
		}
	}

	if deliveryParam.State != nil {
		d.State = string(*deliveryParam.State)
	}

	if deliveryParam.Order.ID != 0 {
		o := new(order.Order)
		if err := db.
			Preload("Items").
			Preload("Vendor").
			Preload("User").
			First(o, "id = ?", deliveryParam.Order.ID).Error; err != nil {
			return nil, cmd.ParseMysqlError(ctx, "order", err)
		}
		d.Order = o
	}

	if len(deliveryParam.DeliveryItems) > 0 {

		orderItems, err := validateOrderItems(ctx, d.Order, deliveryParam.DeliveryItems)

		if err != nil {
			return nil, err
		}

		err = updateOrderStatus(ctx, orderItems, d.Order)

		if err != nil {
			return nil, err
		}

		for _, di := range deliveryParam.DeliveryItems {

			key := orderItems.KeyFunc()(di.ProductID, di.StoreID)

			orderItem, ok := orderItems[key]
			if !ok {
				return nil, net.NewBadRequestError(ctx, errors.New(fmt.Sprintf("item [%d] not found in order", di.ProductID)))
			}

			i := &Item{
				ProductID: di.ProductID,
				StoreID:   di.StoreID,
				Quantity:  di.Quantity,
				UnitPrice: orderItem.UnitPrice,
			}
			d.Items = append(d.Items, i)
		}
	}

	if err := db.Create(d).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "delivery", err)
	}

	l.Info(fmt.Sprintf("Delivery -> created\n [%d]", d.ID))

	return d, nil
}

func updateOrderStatus(ctx context.Context, items map[string]*order.Item, o *order.Order) error {

	l := logging.FromContext(ctx)
	l.Info("checking if order status must be updated")

	status := model.Accepted

	for _, item := range items {
		if item.Quantity != 0 {
			status = model.AcceptedPartially
			l.Debug(fmt.Sprintf("item [%d] -> quantity not zero, order status -> AcceptedPartially", item.ProductID))
			break
		}
	}

	if o.State == status.String() {
		l.Debug(fmt.Sprintf("order status [%s] -> no change", status))
		return nil
	}

	l.Info(fmt.Sprintf("order status [%s] -> updating", status))

	err := order.UpdateOrder(ctx, o.ID, &model.Order{
		State: &status,
	})

	return err
}

type OrderItems map[string]*order.Item

func (oi OrderItems) KeyFunc() func(uint, uint) string {
	return func(productID uint, storeID uint) string {
		return fmt.Sprintf("%d-%d", productID, storeID)
	}
}

func validateOrderItems(ctx context.Context, o *order.Order, deliveryItems []*model.DeliveryItem) (OrderItems, error) {

	orderItemMap := make(OrderItems)
	//order.item tem productId, orderId e storeId como primary key
	fnOrderItemMapKey := orderItemMap.KeyFunc()

	for _, item := range o.Items {
		key := fnOrderItemMapKey(item.ProductID, item.StoreID)
		orderItemMap[key] = item
	}

	//getting other deliveries for given order and subtracting available quantity
	orderDeliveries, err := GetDeliveryByOrderID(ctx, o.ID)

	if err != nil {
		return nil, err
	}

	for _, orderDelivery := range orderDeliveries {
		for _, orderDeliveryItem := range orderDelivery.Items {
			key := fnOrderItemMapKey(orderDeliveryItem.ProductID, orderDeliveryItem.StoreID)
			i, ok := orderItemMap[key]
			if ok {
				i.Quantity -= orderDeliveryItem.Quantity
			}
		}
	}

	for _, di := range deliveryItems {

		key := fnOrderItemMapKey(di.ProductID, di.StoreID)

		itemInOrder, ok := orderItemMap[key]

		if !ok {
			return nil, net.NewHadesError(ctx,
				errors.New(fmt.Sprintf("Product %d not found in order %d", di.ProductID, o.ID)),
				http.StatusBadRequest,
			)
		}

		if di.Quantity > itemInOrder.Quantity {
			return nil, net.NewHadesError(ctx,
				errors.New(fmt.Sprintf("Product %d quantity %f is greater than order quantity %f", di.ProductID, di.Quantity, itemInOrder.Quantity)),
				http.StatusBadRequest,
			)
		}

	}

	return orderItemMap, nil
}

// UpdateDelivery creates a new delivery
func UpdateDelivery(ctx context.Context, deliveryID uint, deliveryParam *model.Delivery) (*Delivery, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	marshal, err := json.Marshal(deliveryParam)

	if err != nil {
		return nil, err
	}

	l.Info(fmt.Sprintf("Updating delivery -> \n [%s]", string(marshal)))

	// verify that the orderParams exists
	existingDelivery := new(Delivery)

	if err := deliveryQuery(db).First(existingDelivery, "id = ?", deliveryID).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "delivery", err)
	}

	if deliveryParam.Session != nil {
		if deliveryParam.Session.ID != 0 {
			s := new(Session)
			if err := db.
				Preload("Vehicle").
				Preload("User").
				First(s, "id = ?", deliveryParam.Session.ID).Error; err != nil {
				return nil, cmd.ParseMysqlError(ctx, "session", err)
			}
			existingDelivery.Session = s
		}
	}

	// for now we don't allow to change the order
	//if deliveryParam.Order.ID != 0 {
	//	o := new(order.Order)
	//	if err := db.Preload("User").First(o, "id = ?", deliveryParam.Order.ID).Error; err != nil {
	//		return nil, cmd.ParseMysqlError(ctx, "order", err)
	//	}
	//	existingDelivery.Order = o
	//}
	tx := db.Begin()

	if len(deliveryParam.DeliveryItems) > 0 {

		orderWithItems := new(order.Order)

		if err := db.
			Preload("Items").
			Preload("Vendor").
			Preload("User").
			First(orderWithItems, "id = ?", existingDelivery.OrderID).Error; err != nil {
			return nil, cmd.ParseMysqlError(ctx, "order", err)
		}

		orderItems, err := validateOrderItems(ctx, orderWithItems, deliveryParam.DeliveryItems)

		if err != nil {
			return nil, err
		}

		err = updateOrderStatus(ctx, orderItems, orderWithItems)

		if err != nil {
			return nil, err
		}

		if err := existingDelivery.updateItems(deliveryParam.DeliveryItems); err != nil {
			return nil, cmd.ParseMysqlError(ctx, "DeliveryItems", err)
		}
	}

	if deliveryParam.State != nil {
		existingDelivery.State = string(*deliveryParam.State)
	}

	if deliveryParam.EndDate != "" {
		parse, err := time.Parse(time.RFC3339, deliveryParam.EndDate)

		if err != nil {
			return nil, net.NewHadesError(ctx, err, http.StatusBadRequest)
		}

		existingDelivery.EndDate.Time = parse
		existingDelivery.EndDate.Valid = true
	}

	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(existingDelivery).Error; err != nil {
		tx.Rollback()
		return nil, cmd.ParseMysqlError(ctx, "delivery", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, cmd.ParseMysqlError(ctx, "delivery", err)
	}

	l.Info(fmt.Sprintf("Delivery -> updated\n [%d]", existingDelivery.ID))

	return existingDelivery, nil
}

// updateItems updates the items of the order - NO ITEMS ARE REMOVED
func (d *Delivery) updateItems(newItems []*model.DeliveryItem) error {

	fnKey := func(productId, storeId uint) string {
		return fmt.Sprintf("%d-%d", productId, storeId)
	}

	mapItem := make(map[string]*Item)

	for _, i := range d.Items {
		mapItem[fnKey(i.ProductID, i.StoreID)] = i
	}

	for _, i := range newItems {
		key := fnKey(i.ProductID, i.StoreID)
		if _, ok := mapItem[key]; ok {
			mapItem[key].Quantity = i.Quantity
		} else {
			mapItem[key] = &Item{
				ProductID: i.ProductID,
				StoreID:   i.StoreID,
				Quantity:  i.Quantity,
			}
		}
	}

	d.Items = make([]*Item, 0)

	for _, i := range mapItem {
		d.Items = append(d.Items, i)
	}

	return nil

}

func GetDeliveries(ctx context.Context, opts *GetDeliveryOptions) ([]*Delivery, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Getting deliveries -> \n [%s]", opts.Params.Encode()))

	var deliveries []*Delivery

	query := deliveryQuery(db)
	query = opts.parseDeliveryParams(query)

	if err := query.Find(&deliveries).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "delivery", err)
	}

	return deliveries, nil
}

func GetDelivery(ctx context.Context, deliveryID uint) (*Delivery, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Getting delivery -> \n [%d]", deliveryID))

	delivery := new(Delivery)

	if err := deliveryQuery(db).First(delivery, "id = ?", deliveryID).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "delivery", err)
	}

	return delivery, nil
}

func GetDeliveryByOrderID(ctx context.Context, orderID uint) ([]*Delivery, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Getting delivery by orderID -> \n [%d]", orderID))

	var deliveries []*Delivery

	if err := db.
		Preload("Items").
		Preload("Items.Product").
		Find(&deliveries, "order_id = ?", orderID).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "delivery", err)
	}

	return deliveries, nil
}

func DeleteDelivery(ctx context.Context, deliveryID uint) error {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Deleting delivery -> \n [%d]", deliveryID))

	delivery := new(Delivery)

	if err := db.First(delivery, "id = ?", deliveryID).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "delivery", err)
	}

	if err := db.Delete(delivery).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "delivery", err)
	}

	l.Info(fmt.Sprintf("Deleted delivery [%d]", deliveryID))

	return nil
}

func RemoveDeliveryItems(ctx context.Context, deliveryID uint, items []*model.DeliveryItem) error {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Removing delivery items -> \n [%d]", deliveryID))

	delivery := new(Delivery)

	if err := db.First(delivery, "id = ?", deliveryID).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "delivery", err)
	}

	itemsToRemove := make([]*Item, 0)

	for _, item := range items {
		i := &Item{ProductID: item.ProductID, StoreID: item.StoreID, DeliveryID: deliveryID}
		itemsToRemove = append(itemsToRemove, i)
	}

	if err := db.Delete(itemsToRemove).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "delivery_items", err)
	}

	l.Info(fmt.Sprintf("Removed delivery items [%d]", deliveryID))

	return nil
}

func deliveryQuery(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Store").
		Preload("Session").
		Preload("Session.Vehicle").
		Preload("Session.User").
		Preload("Order")
}

// GetDeliveryOptions TODO: add pagination
type GetDeliveryOptions struct {
	Params url.Values
}

func (o *GetDeliveryOptions) parseDeliveryParams(query *gorm.DB) *gorm.DB {

	tableName := (&Delivery{}).TableName()

	if s := o.Params.Get("state"); s != "" {
		query = query.Where(tableName+".state = ?", s)
	}

	if s := o.Params.Get("user_id"); s != "" {
		query = query.
			Joins("INNER JOIN sessions ON sessions.id = deliveries.session_id").
			Where("sessions.user_id = ?", s)
	}

	if s := o.Params.Get("vehicle_id"); s != "" {
		query = query.
			Joins("INNER JOIN sessions ON sessions.id = deliveries.session_id").
			Where("sessions.vehicle_id = ?", s)
	}

	if s := o.Params.Get("order_id"); s != "" {
		query = query.Where(tableName+".order_id = ?", s)
	}

	if s := o.Params.Get("created_at"); s != "" {
		//TODO: parse date
		query = query.Where(tableName+".created_at >= ?", s)
	}

	if s := o.Params.Get("end_date"); s != "" {
		//TODO: parse date
		query = query.Where(tableName+".end_date <= ?", s)
	}

	return query
}
