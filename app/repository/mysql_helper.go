package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"hades_backend/api/utils/net"
	"net/http"
)

func ParseMysqlError(ctx context.Context, entity string, err error) error {

	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return net.NewHadesError(ctx, errors.New(fmt.Sprintf("%s not found", entity)), http.StatusNotFound)
	}

	var mysqlErr *mysql.MySQLError

	isMysqlError := errors.As(err, &mysqlErr)

	if isMysqlError {
		switch mysqlErr.Number {
		case 1062:
			return net.NewHadesError(ctx, errors.New(fmt.Sprintf("%s already exists", entity)), http.StatusConflict)
		case 1452:
			return net.NewHadesError(ctx, mysqlErr, http.StatusBadRequest)
		}
	}

	return net.NewHadesError(ctx, err, http.StatusInternalServerError)
}
