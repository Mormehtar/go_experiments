package interfaces

type IName interface {
	Create(name string) (*Name, error)
	Get(name string) (*Name, error)
	Update(id int64, name string) (*Name, error)
	Delete(id int64) error
}

type IProperty interface {
	Create(nameId int64, key string, value string) (*Property, error)
	Get(nameId int64, key string) ([]*Property, error)
	Update(id int64, key string, value string) (*Property, error)
	Delete(id int64) error
}
