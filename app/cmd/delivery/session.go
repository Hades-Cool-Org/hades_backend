package delivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd"
	"hades_backend/app/cmd/user"
	"hades_backend/app/database"
	"hades_backend/app/logging"
	"hades_backend/app/model"
	"net/url"
	"time"
)

type Session struct {
	gorm.Model

	UserID uint       `gorm:"index"`
	User   *user.User //motorista

	VehicleID uint `gorm:"index"`
	Vehicle   *Vehicle

	EndDate sql.NullTime `gorm:"index"`
}

func (s Session) TableName() string {
	return "sessions"
}

func (s *Session) BeforeDelete(tx *gorm.DB) error {
	if !s.EndDate.Valid {
		loc, _ := time.LoadLocation("America/Sao_Paulo")

		x := time.Now().In(loc)

		s.EndDate.Time = x
		s.EndDate.Valid = true
	}
	return nil
}

func CreateSession(ctx context.Context, sessionParam *model.Session) (*Session, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Creating session -> \n [%+v]", sessionParam))

	s := new(Session)

	/// checking for current sessions for the user/vehicle
	if err := db.
		First(s, "(user_id = ? OR vehicle_id = ?)",
			sessionParam.User.ID,
			sessionParam.Vehicle.ID,
		).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.Info("Session not found, creating new one")
		} else {
			return nil, cmd.ParseMysqlError(ctx, "session", err)
		}
	}

	if s != nil {
		return nil, net.NewHadesError(ctx, errors.New("session already exists"), 400)
	}

	s.UserID = sessionParam.User.ID
	s.VehicleID = sessionParam.Vehicle.ID

	if err := database.DB.Create(s).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "session", err)
	}

	return s, nil
}

func GetSession(ctx context.Context, sessionID uint) (*Session, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Getting session -> \n [%d]", sessionID))

	s := new(Session)

	if err := db.
		Preload("Vehicle").
		Preload("User").
		First(s, "id = ?", sessionID).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "session", err)
	}

	return s, nil
}

func GetSessions(ctx context.Context, options *GetSessionOptions) ([]*Session, error) {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Getting sessions -> \n [%+v]", options))

	var sessions []*Session

	query := db.
		Preload("Vehicle").
		Preload("User").
		Order("created_at DESC")

	if options != nil {
		query = options.parseSessionParam(query)
	}

	if err := query.Find(&sessions).Error; err != nil {
		return nil, cmd.ParseMysqlError(ctx, "session", err)
	}

	return sessions, nil
}

func DeleteSession(ctx context.Context, sessionID uint) error {

	l := logging.FromContext(ctx)
	db := database.DB.WithContext(ctx)

	l.Info(fmt.Sprintf("Deleting session -> \n [%d]", sessionID))

	s := new(Session)

	if err := db.First(s, "id = ?", sessionID).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "session", err)
	}

	if err := database.DB.Delete(s).Error; err != nil {
		return cmd.ParseMysqlError(ctx, "session", err)
	}

	return nil
}

// GetSessionOptions TODO: add pagination
type GetSessionOptions struct {
	Params url.Values
}

func (o *GetSessionOptions) parseSessionParam(query *gorm.DB) *gorm.DB {

	tableName := (&Session{}).TableName()

	if s := o.Params.Get("vehicle_id"); s != "" {
		query = query.Where(tableName+".vehicle_id = ?", s)
	}

	if s := o.Params.Get("user_id"); s != "" {
		query = query.Where(tableName+".user_id = ?", s)
	}

	if s := o.Params.Get("active"); s != "" {
		if s == "true" {
			query = query.Where(tableName + ".deleted_at IS NULL")
		} else {
			query = query.Where(tableName + ".deleted_at IS NOT NULL")
		}
	}

	return query
}
