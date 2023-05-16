package conference

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hades_backend/app/cmd"
	"hades_backend/app/cmd/delivery"
	"hades_backend/app/cmd/occurence"
	"hades_backend/app/cmd/stock"
	"hades_backend/app/database"
	"hades_backend/app/logging"
	"hades_backend/app/model"
)

func DoConference(ctx context.Context, params *model.Occurrence) error {

	dbz := database.DB.WithContext(ctx)
	l := logging.FromContext(ctx)

	tx := dbz.Begin()

	_, err := occurence.CreateOccurrence(ctx, tx, params)

	if err != nil {
		tx.Rollback()
		return err
	}

	storeID := params.StoreID

	l.Info(fmt.Sprintf("Getting stock for store %d", storeID))

	s := &stock.Stock{}

	if result := tx.
		First(s, "store_id = ?", storeID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// should we really do that?!
			s, err = stock.CreateStock(ctx, tx, &model.Stock{
				Store: &model.Store{ID: storeID},
			})
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return result.Error
		}
	}

	del := new(delivery.Delivery)
	if err := tx.
		Preload("Items").
		Preload("Order").
		First(del, "id = ?", params.DeliveryID).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "delivery", err)
	}

	deliveryItemsMap := make(map[string]*delivery.Item)
	fnKey := func(productID uint, storeID uint) string {
		return fmt.Sprintf("%d-%d", productID, storeID)
	}

	for _, item := range del.Items {
		deliveryItemsMap[fnKey(item.ProductID, item.StoreID)] = item
	}

	stockItems := make([]*model.StockItem, 0)

	for _, item := range params.Items {

		key := fnKey(item.ProductID, storeID)

		deliveryItem, ok := deliveryItemsMap[key]

		if !ok {
			tx.Rollback()
			return errors.New("delivery item not found")
		}

		stockItem := &model.StockItem{
			ProductID: item.ProductID,
			Current:   item.Quantity,
			AvgPrice:  deliveryItem.UnitPrice,
		}

		stockItems = append(stockItems, stockItem)
	}

	err = stock.AddStockItem(ctx, tx, s.ID, stockItems)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return cmd.ParseMysqlError(ctx, "conference", err)
	}

	return nil
}
