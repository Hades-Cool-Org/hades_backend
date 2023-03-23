package user

import (
	"context"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hades_backend/app/model"
	"regexp"
	"testing"
)

type AnyDBValue struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyDBValue) Match(v driver.Value) bool {
	return true
}

func TestMySqlRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	open, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	m := &MySqlRepository{db: open}

	u := &model.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "1234",
	}

	expectedID := uint(1)
	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`"+
			",`name`,`email`,`phone`,`password`,`first_login`) VALUES (?,?,?,?,?,?,?,?)")).
		WithArgs(AnyDBValue{}, AnyDBValue{}, AnyDBValue{}, u.Name, u.Email, u.Phone, AnyDBValue{}, true).
		WillReturnResult(sqlmock.NewResult(int64(expectedID), 1))

	mock.ExpectCommit()

	id, err := m.Create(context.TODO(), u)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if id != expectedID {
		t.Errorf("expected ID to be %d, but got %d", expectedID, id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

//func TestMySqlRepository_Update(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer db.Close()
//
//	open, err := gorm.Open(mysql.New(mysql.Config{
//		Conn:                      db,
//		SkipInitializeWithVersion: true,
//	}), &gorm.Config{})
//
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	m := &MySqlRepository{db: open}
//
//	u := &user.User{
//		ID:    uint(1),
//		Name:  "John Doe",
//		Email: "john.doe@example.com",
//	}
//
//	mock.ExpectBegin()
//	mock.ExpectExec("DELETE FROM `roles`").WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))
//
//	mock.ExpectExec(
//		regexp.QuoteMeta("UPDATE `users` SET `updated_at`=?,`name`=?,`email`=?"+
//			" WHERE `users`.`deleted_at` IS NULL AND `id` = ?")).
//		WithArgs(AnyDBValue{}, u.Name, u.Email, u.ID).
//		WillReturnResult(sqlmock.NewResult(0, 1))
//
//	mock.ExpectCommit()
//
//	err = m.Update(context.TODO(), u)
//	if err != nil {
//		t.Errorf("unexpected error: %v", err)
//	}
//
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %v", err)
//	}
//}
