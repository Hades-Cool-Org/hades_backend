package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"hades_backend/api/utils/net"
	"net/http"
)

func ParseMysqlError(ctx context.Context, entity string, err error) error {

	if err == nil {
		return nil
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
