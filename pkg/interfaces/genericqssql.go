package interfaces

type GenericQSSQL interface {
	GetQueryString() string
	GetParsedQueryString() map[string]interface{}
	GetTemplate() string

	AddWhereAlias(string, string)
	AddWhereAliases(map[string]string)
	GetWhereAliases() map[string]string

	AddWhereValue(string)
	AddWhereValues([]string)
	GetWhereValues() []string
	AddWhereBindValue(interface{})
	GetWhereBindValues() []interface{}
	GetWhereClause() string

	Parse() error
	GetSQL() (string, error)
	GetBindValues() []interface{}

	SetParamLimitValue(uint64)
	GetParamLimitName() string
	GetLimitClause() string

	SetParamOffsetValue(uint64)
	GetParamOffsetName() string
	GetOffsetClause() string

	AddParamOrderByValue(string)
	GetParamOrderByName() string
	GetParamOrderByValues() []string
	GetOrderByClause() string
}
