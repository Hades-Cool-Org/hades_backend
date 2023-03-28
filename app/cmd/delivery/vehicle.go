package delivery

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"hades_backend/app/cmd"
	"hades_backend/app/database"
	"hades_backend/app/logging"
	"hades_backend/app/model"
)

type Vehicle struct {
	gorm.Model
	Name string
	Type string
}

func (v Vehicle) TableName() string {
	return "vehicles"
}

func CreateVehicle(ctx context.Context, vehicleParam *model.Vehicle) (*Vehicle, error) {
	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Creating vehicle -> \n [%+v]", vehicleParam))

	v := new(Vehicle)
	v.Name = vehicleParam.Name
	v.Type = vehicleParam.Type

	if err := db.Create(v).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "vehicle", err)
	}

	return v, nil
}

func GetVehicle(ctx context.Context, id uint) (*Vehicle, error) {
	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Getting vehicle -> \n [%d]", id))

	v := new(Vehicle)

	if err := db.First(v, id).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "vehicle", err)
	}

	return v, nil
}

func UpdateVehicle(ctx context.Context, id uint, vehicleParam *model.Vehicle) (*Vehicle, error) {
	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Updating vehicle -> \n [%d]", id))

	v := new(Vehicle)

	if err := db.First(v, id).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "vehicle", err)
	}

	v.Name = vehicleParam.Name
	v.Type = vehicleParam.Type

	if err := db.Save(v).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "vehicle", err)
	}

	return v, nil
}

func DeleteVehicle(ctx context.Context, id uint) error {
	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Deleting vehicle -> \n [%d]", id))

	v := new(Vehicle)

	if err := db.First(v, id).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "vehicle", err)
	}

	if err := db.Delete(v).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "vehicle", err)
	}

	return nil
}

func GetAllVehicles(ctx context.Context) ([]*Vehicle, error) {
	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Listing vehicles"))

	v := make([]*Vehicle, 0)

	if err := db.Find(&v).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "vehicle", err)
	}

	return v, nil
}
