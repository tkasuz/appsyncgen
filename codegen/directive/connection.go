package directive

type ConnectionType string

const (
	HasOne     ConnectionType = "hasOne"
	HasMany    ConnectionType = "hasMany"
	ManyToMany ConnectionType = "manyToMany"
)

func (c ConnectionType) IsHasOne() bool {
	if c == HasOne {
		return true
	} else {
		return false
	}
}

func (c ConnectionType) IsHasMany() bool {
	if c == HasMany {
		return true
	} else {
		return false
	}
}

func (c ConnectionType) IsManyToMany() bool {
	if c == ManyToMany {
		return true
	} else {
		return false
	}
}

func NewConnection(name string) *ConnectionType {
	if IsConnectionType(name) {
		con := ConnectionType(name)
		return &con
	}
	return nil
}

func IsConnectionType(v string) bool {
	return v == string(HasOne) || v == string(HasMany) || v == string(ManyToMany)
}
