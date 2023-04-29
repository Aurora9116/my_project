package database

type DbConn interface {
	RollBack()
	Commit()
	Begin()
}
