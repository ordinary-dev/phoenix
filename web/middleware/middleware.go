package middleware

import (
	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web/controllers"
)

type Middleware struct {
	ctrl *controllers.Controllers
	db   database.Database
}

func New(ctrl *controllers.Controllers, db database.Database) Middleware {
	return Middleware{
		ctrl: ctrl,
		db:   db,
	}
}
