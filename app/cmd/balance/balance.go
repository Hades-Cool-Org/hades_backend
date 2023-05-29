package balance

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"hades_backend/app/cmd"
	"hades_backend/app/cmd/user"
	"hades_backend/app/database"
	"net/http"
	"sync"
)

type Balance struct {
	gorm.Model

	UserID uint `gorm:"index"`
	User   *user.User

	Amount decimal.Decimal `gorm:"type:decimal(12,3);"`

	Entries []*Entry

	ModificationLock sync.Mutex `json:"-" sql:"-" gorm:"-"`
}

func (b *Balance) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Entry struct {
	gorm.Model

	BalanceID uint `gorm:"index"`

	Amount decimal.Decimal `gorm:"type:decimal(12,3);"`

	// 0 = credit, 1 = debit
	Type int
}

func (b Entry) TableName() string {
	return "balance_entry"
}

func GetBalance(ctx context.Context, db *gorm.DB, userID uint) (*Balance, error) {

	b := new(Balance)

	if err := db.First(b, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, cmd.ParseMysqlError(ctx, "balance", err)
		}
	}

	return b, nil
}

type Operation int

const (
	Credit Operation = iota
	Debit
)

type Params struct {
	UserID    uint
	Amount    decimal.Decimal
	Operation Operation
}

func (p *Params) GetAmount() decimal.Decimal {
	if p.Operation == Credit {
		return p.Amount
	}

	return p.Amount.Neg()
}

func CreateBalance(ctx context.Context, db *gorm.DB, userID uint) (*Balance, error) {

	b := new(Balance)
	b.UserID = userID
	b.Amount = decimal.Zero

	if err := db.Create(b).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "balance", err)
	}

	return b, nil
}

func ManageBalance(ctx context.Context, params *Params) (*Balance, error) {

	dbz := database.DB.WithContext(ctx)
	tx := dbz.Begin()
	b := new(Balance)

	if result := tx.
		First(b, "user_id = ?", params.UserID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// should we really do that?!
			balance, err := CreateBalance(ctx, tx, params.UserID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			b = balance
		} else {
			tx.Rollback()
			return nil, result.Error
		}
	}

	b.ModificationLock.Lock()
	defer b.ModificationLock.Unlock()

	b.Amount = b.Amount.Add(params.GetAmount())

	if err := tx.Save(b).Error; err != nil {
		tx.Rollback()
		return nil, cmd.ParseMysqlError(ctx, "balance", err)
	}

	return b, nil
}
