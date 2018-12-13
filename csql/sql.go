package csql

import (
	"context"
	"fmt"

	"github.com/tusharsoni/copper/clogger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    Config
	Logger    clogger.Logger
}

func NewGormDB(p Params) (*gorm.DB, error) {
	conn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		p.Config.Host,
		p.Config.Port,
		p.Config.User,
		p.Config.Name,
	)
	p.Logger.Infow("Connecting to database..", map[string]string{
		"connection": conn,
	})

	if p.Config.Password != "" {
		conn = fmt.Sprintf("%s password=%s", conn, p.Config.Password)
	}

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		p.Logger.Errorw("Failed to connect to database..", map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	p.Lifecycle.Append(fx.Hook{
		OnStop: func(context context.Context) error {
			p.Logger.Infow("Closing connection to database..", nil)
			return db.Close()
		},
	})

	return db, nil
}
