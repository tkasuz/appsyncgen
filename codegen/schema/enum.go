package schema

type Enum struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type EnumList []*Enum

func (l EnumList) IsEnum(v string) bool {
	for _, it := range l {
		if it.Name == v {
			return true
		}
	}
	return false
}

func (l EnumList) ForName(name string) *Enum {
	for _, it := range l {
		if it.Name == name {
			return it
		}
	}
	return nil
}
