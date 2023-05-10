package conference

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hades_backend/app/cmd"
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
		//Preload("items").
		//Preload("Items.Product").
		First(s, "store_id = ?", storeID); result.Error != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
			return err
		}
	}

	stockItems := make([]*model.StockItem, len(params.Items))

	for _, item := range params.Items {

		stockItem := &model.StockItem{
			ProductID: item.ProductID,
			Current:   item.Quantity,
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
