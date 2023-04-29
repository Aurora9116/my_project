package tran

import "test.com/project-project/internal/database"

type Transaction interface {
	Action(func(conn database.DbConn) error) error
}
