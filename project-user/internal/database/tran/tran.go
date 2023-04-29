package tran

import "test.com/project-user/internal/database"

type Transaction interface {
	Action(func(conn database.DbConn) error) error
}
