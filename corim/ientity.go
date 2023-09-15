package corim

type IEntity interface {
	SetEntityName(name string) IEntity
	SetRegID(uri string) IEntity
	SetRoles(roles ...Role) IEntity
	Valid() error
	FromCBOR(data []byte) error
	ToCBOR() ([]byte, error)
}
