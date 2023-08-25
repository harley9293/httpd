package orm

type ORM interface {
	Connect() error
	Migrate(...any) error
}
