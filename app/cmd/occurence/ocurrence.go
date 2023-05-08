package occurence

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd"
	"hades_backend/app/cmd/delivery"
	"hades_backend/app/cmd/store"
	"hades_backend/app/cmd/user"
	"hades_backend/app/database"
	"hades_backend/app/logging"
	"hades_backend/app/model"
	"net/url"
)

type Occurrence struct {
	gorm.Model

	OrderID  uint `gorm:"index;not null;"`
	VendorID uint `gorm:"index;not null;"`

	DeliveryID uint `gorm:"not null;"`
	Delivery   *delivery.Delivery

	StoreID uint `gorm:"not null;"`
	Store   *store.Store

	UserID uint
	User   *user.User

	Items []*Item
}

type Item struct {
	OccurrenceID uint `gorm:"primaryKey;autoIncrement:false"`

	ProductID     uint `gorm:"primaryKey;autoIncrement:false"`
	Name          string
	MeasuringUnit string

	Type     string // CREDIT/DEBIT
	Quantity float64

	UnitPrice decimal.Decimal `gorm:"type:decimal(12,3);"`
}

func (Item) TableName() string {
	return "occurrence_items"
}

func (Occurrence) TableName() string {
	return "occurrences"
}

func (o *Occurrence) BeforeDelete(tx *gorm.DB) error {

	delModels := map[string]interface{}{
		"items": &[]Item{},
	}
	for name, dm := range delModels {
		if result := tx.Delete(dm, "occurrence_id = ?", o.ID); result.Error != nil {
			return errors.Wrap(result.Error, fmt.Sprintf("Error deleting %s records", name))
		}
	}
	return nil
}

func CreateOccurrence(ctx context.Context, occurrenceParams *model.Occurrence) (*Occurrence, error) {
	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	marshal, err := json.Marshal(occurrenceParams)

	if err != nil {
		return nil, err
	}

	l.Info(fmt.Sprintf("Creating occurrence -> \n [%s]", string(marshal)))

	o := &Occurrence{}

	if occurrenceParams.User.ID != 0 {
		u := new(user.User)
		if err := db.First(u, "id = ?", occurrenceParams.User.ID).Error; err != nil {
			return nil, cmd.ParseMysqlError(ctx, "user", err)
		}
		o.User = u
	}

	if occurrenceParams.DeliveryID != 0 {
		del := new(delivery.Delivery)
		if err := db.
			Preload("Items").
			Preload("Order").
			First(del, "id = ?", occurrenceParams.DeliveryID).Error; err != nil {
			return nil, cmd.ParseMysqlError(ctx, "delivery", err)
		}
		o.Delivery = del
		o.DeliveryID = del.ID
		o.OrderID = del.OrderID
	} else {
		return nil, net.NewBadRequestError(ctx, errors.New("delivery id is required"))
	}

	if occurrenceParams.StoreID != 0 {
		st := new(store.Store)
		if err := db.
			First(st, "id = ?", occurrenceParams.StoreID).Error; err != nil {
			return nil, cmd.ParseMysqlError(ctx, "store", err)
		}
		o.Store = st
		o.StoreID = st.ID
	} else {
		return nil, net.NewBadRequestError(ctx, errors.New("store id is required"))
	}

	o.Items = generateItems(occurrenceParams, o.Delivery.Items)
	o.VendorID = o.Delivery.Order.VendorID

	if err := db.Create(o).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "occurrence", err)
	}

	return o, nil
}

// generateItems logic to generate debit/credit items
func generateItems(occurrenceParams *model.Occurrence, deliveryItems []*delivery.Item) []*Item {
	storeID := occurrenceParams.StoreID

	storeItemsInOrder := make(map[uint]*delivery.Item)

	for _, item := range deliveryItems {
		if item.StoreID == storeID {
			storeItemsInOrder[item.ProductID] = item
		}
	}

	items := make([]*Item, len(occurrenceParams.Items))

	for _, requestItem := range occurrenceParams.Items {

		deliveryItem := storeItemsInOrder[requestItem.ProductID]

		if deliveryItem == nil {
			// caso recebermos item que não está no pedido, lojão tem um debito com o fornecedor
			items = append(items, &Item{
				ProductID: requestItem.ProductID,
				Type:      "DEBIT", //todo enum?
				Quantity:  requestItem.Quantity,
				UnitPrice: deliveryItem.UnitPrice,
			})
			continue
		}

		if requestItem.Quantity == deliveryItem.Quantity {
			// caso recebermos um item com a mesma quantidade do delivery, lojão não tem debito nem credito com o fornecedor
			continue
		}

		if requestItem.Quantity > deliveryItem.Quantity {
			// caso recebermos um item com uma quantidade maior do que o delivery, lojão tem um debito com o fornecedor
			quantity := requestItem.Quantity - deliveryItem.Quantity

			items = append(items, &Item{
				ProductID: requestItem.ProductID,
				Type:      "DEBIT", //todo enum?
				Quantity:  quantity,
			})
			continue
		}

		if requestItem.Quantity < deliveryItem.Quantity {
			// caso recebermos um item com uma quantidade menor do que o delivery, lojão tem um credito com o fornecedor
			quantity := deliveryItem.Quantity - requestItem.Quantity

			items = append(items, &Item{
				ProductID: requestItem.ProductID,
				Type:      "CREDIT", //todo enum?
				Quantity:  quantity,
			})
			continue
		}
	}

	return items
}

func DeleteOccurrence(ctx context.Context, occurrenceID uint) error {
	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	o := new(Occurrence)

	if err := db.First(o, "id = ?", occurrenceID).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "occurrence", err)
	}

	l.Info(fmt.Sprintf("Deleting occurrence -> \n [%+v]", o))

	if err := db.Delete(o).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "occurrence", err)
	}

	return nil
}

func GetOccurrence(ctx context.Context, occurrenceID uint) (*Occurrence, error) {
	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	o := new(Occurrence)

	if err := db.
		Preload("Items").
		First(o, "id = ?", occurrenceID).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "occurrence", err)
	}

	l.Info(fmt.Sprintf("Get occurrence -> \n [%+v]", o))

	return o, nil
}

func GetOccurrences(ctx context.Context, options *GetOccurrenceOptions) ([]*Occurrence, error) {
	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	query := db.Preload("Items")

	query = options.parseDeliveryParams(query)

	occurrences := make([]*Occurrence, 0)

	if err := query.Find(&occurrences).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "occurrence", err)
	}

	l.Info(fmt.Sprintf("Get occurrences -> \n [%+v]", occurrences))

	return occurrences, nil
}

// GetOccurrenceOptions TODO: add pagination
type GetOccurrenceOptions struct {
	Params url.Values
}

func (o *GetOccurrenceOptions) parseDeliveryParams(query *gorm.DB) *gorm.DB {
	if o.Params.Get("delivery_id") != "" {
		query = query.Where("delivery_id = ?", o.Params.Get("delivery_id"))
	}

	if o.Params.Get("order_id") != "" {
		query = query.Where("order_id = ?", o.Params.Get("order_id"))
	}

	if o.Params.Get("store_id") != "" {
		query = query.Where("store_id = ?", o.Params.Get("store_id"))
	}

	if o.Params.Get("vendor_id") != "" {
		query = query.Where("vendor_id = ?", o.Params.Get("vendor_id"))
	}

	//todo fetch deleted?

	return query
}
