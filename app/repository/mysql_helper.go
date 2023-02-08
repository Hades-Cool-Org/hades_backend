package repository

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"hades_backend/app/hades_errors"
	"net/http"
)

func ParseMysqlError(entity string, err error) error {

	var mysqlErr *mysql.MySQLError

	isMysqlError := errors.As(err, &mysqlErr)

	if isMysqlError {

		switch mysqlErr.Number {
		case 1062:
			return hades_errors.NewHadesError(errors.New(fmt.Sprintf("%s already exists", entity)), http.StatusConflict)
		}
	}

	return hades_errors.NewHadesError(err, http.StatusInternalServerError)
}
