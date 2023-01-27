package datasource

type DataSource struct {
	Name string         `json:"name"`
	Type DataSourceType `json:"type"`
}

type DataSourceList []*DataSource

func (l DataSourceList) ForName(name string) *DataSource {
	for _, it := range l {
		if it.Name == name {
			return it
		}
	}
	return nil
}

type DataSourceType string

const (
	NONE     DataSourceType = "NONE"
	DYNAMODB DataSourceType = "DYNAMODB"
)

func NewDataSource(datasourceType string, datasourceName string) *DataSource {
	if datasourceType == string(NONE) || datasourceType == string(DYNAMODB) {
		switch datasourceType {
		case string(DYNAMODB):
			return &DataSource{
				Name: datasourceName,
				Type: DYNAMODB,
			}
		}
	}
	return nil
}
