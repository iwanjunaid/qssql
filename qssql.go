package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	helperHaving "github.com/iwanjunaid/qssql/internal/helpers/having"
	helperLimit "github.com/iwanjunaid/qssql/internal/helpers/limit"
	helperOffset "github.com/iwanjunaid/qssql/internal/helpers/offset"
	helperOrderBy "github.com/iwanjunaid/qssql/internal/helpers/orderby"
	helperWhere "github.com/iwanjunaid/qssql/internal/helpers/where"
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
	"github.com/iwanjunaid/qssql/pkg/parser/qs"
)

const (
	DEFAULT_PARAM_LIMIT_NAME  = "limit"
	DEFAULT_PARAM_LIMIT_VALUE = 10

	DEFAULT_PARAM_OFFSET_NAME  = "offset"
	DEFAULT_PARAM_OFFSET_VALUE = 0

	DEFAULT_PARAM_ORDER_BY_NAME = "sort"

	WHERE_KEYWORD    = "WHERE"
	HAVING_KEYWORD   = "HAVING"
	LIMIT_KEYWORD    = "LIMIT"
	OFFSET_KEYWORD   = "OFFSET"
	ORDER_BY_KEYWORD = "ORDER BY"

	AND_OPERATOR = "AND"
)

type QSSQL struct {
	queryString       string
	parsedQueryString map[string]interface{}
	template          string
	sql               string

	beforeHook func(iface.GenericQSSQL)

	whereValues     []string
	whereBindValues []interface{}
	whereAliases    map[string]string

	havingValues     []string
	havingBindValues []interface{}
	havingAliases    map[string]string

	paramLimitName  string
	paramLimitValue uint64
	limitClause     string

	paramOffsetName  string
	paramOffsetValue uint64
	offsetClause     string

	paramOrderByName   string
	paramOrderByValues []string
	orderByClause      string
}

func New(qs, template string, beforeHook func(qssql iface.GenericQSSQL)) iface.GenericQSSQL {
	whereAliases := make(map[string]string)
	havingAliases := make(map[string]string)

	qssql := &QSSQL{
		queryString:      qs,
		template:         template,
		beforeHook:       beforeHook,
		whereAliases:     whereAliases,
		havingAliases:    havingAliases,
		paramLimitName:   DEFAULT_PARAM_LIMIT_NAME,
		paramLimitValue:  DEFAULT_PARAM_LIMIT_VALUE,
		paramOffsetName:  DEFAULT_PARAM_OFFSET_NAME,
		paramOffsetValue: DEFAULT_PARAM_OFFSET_VALUE,
		paramOrderByName: DEFAULT_PARAM_ORDER_BY_NAME,
	}

	return qssql
}

func (q *QSSQL) AddWhereValue(value string) {
	q.whereValues = append(q.whereValues, value)
}

func (q *QSSQL) AddWhereValues(values []string) {
	for _, v := range values {
		q.whereValues = append(q.whereValues, v)
	}
}

func (q *QSSQL) GetWhereValues() []string {
	return q.whereValues
}

func (q *QSSQL) GetWhereClause() string {
	if len(q.whereValues) == 0 {
		return ""
	}

	return WHERE_KEYWORD + " " + strings.Join(q.whereValues[:], " AND ")
}

func (q *QSSQL) AddWhereAlias(key string, value string) {
	q.whereAliases[key] = value
}

func (q *QSSQL) AddWhereAliases(aliases map[string]string) {
	for key, value := range aliases {
		q.whereAliases[key] = value
	}
}

func (q *QSSQL) GetWhereAliases() map[string]string {
	return q.whereAliases
}

func (q *QSSQL) SetParamLimitValue(val uint64) {
	q.paramLimitValue = val
}

func (q *QSSQL) GetParamLimitName() string {
	return q.paramLimitName
}

func (q *QSSQL) SetParamOffsetValue(val uint64) {
	q.paramOffsetValue = val
}

func (q *QSSQL) GetParamOffsetName() string {
	return q.paramOffsetName
}

func (q *QSSQL) AddParamOrderByValue(value string) {
	q.paramOrderByValues = append(q.paramOrderByValues, value)
}

func (q *QSSQL) GetParamOrderByValues() []string {
	return q.paramOrderByValues
}

func (q *QSSQL) GetParamOrderByName() string {
	return q.paramOrderByName
}

func (q *QSSQL) GetQueryString() string {
	return q.queryString
}

func (q *QSSQL) GetParsedQueryString() map[string]interface{} {
	return q.parsedQueryString
}

func (q *QSSQL) GetTemplate() string {
	return q.template
}

func (q *QSSQL) Parse() error {
	parsed, parseErr := qs.ToJSON(q.queryString)

	if parseErr != nil {
		return parseErr
	}

	var data map[string]interface{}

	unmarshallErr := json.Unmarshal([]byte(parsed), &data)

	if unmarshallErr != nil {
		return unmarshallErr
	}

	q.parsedQueryString = data

	// Call before hook
	if q.beforeHook != nil {
		q.beforeHook(q)
	}

	for _key, value := range data {
		keyTrimmed := strings.TrimSpace(_key)
		key := strings.ToLower(keyTrimmed)

		if key == q.paramLimitName || key == q.paramOffsetName ||
			key == q.paramOrderByName {
			if key == q.paramLimitName {
				err := helperLimit.Extract(q, key, value)

				if err != nil {
					return err
				}
			}

			if key == q.paramOffsetName {
				err := helperOffset.Extract(q, key, value)

				if err != nil {
					return err
				}
			}

			if key == q.paramOrderByName {
				err := helperOrderBy.Extract(q, key, value)

				if err != nil {
					return err
				}
			}

			continue
		}

		switch value.(type) {
		case string:
			err := helperWhere.ExtractString(q, key, value)

			if err != nil {
				return err
			}
		case float64:
			err := helperWhere.ExtractFloat64(q, key, value)

			if err != nil {
				return err
			}
		case []interface{}:
			err := helperWhere.ExtractIn(q, key, value)

			if err != nil {
				return err
			}
		case map[string]interface{}:
			realValue := value.(map[string]interface{})

			for _rk, _rv := range realValue {
				rkTrimmed := strings.TrimSpace(_rk)
				rk := strings.ToLower(rkTrimmed)

				if rk == "like" {
					err := helperWhere.ExtractLike(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "hlike" {
					err := helperHaving.ExtractLike(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "ilike" {
					err := helperWhere.ExtractILike(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "hilike" {
					err := helperHaving.ExtractILike(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "ne" {
					err := helperWhere.ExtractNotEqual(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "hne" {
					err := helperHaving.ExtractNotEqual(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "in" {
					err := helperWhere.ExtractIn(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "hin" {
					err := helperHaving.ExtractIn(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "nin" {
					err := helperWhere.ExtractNotIn(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "hnin" {
					err := helperHaving.ExtractNotIn(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "gt" {
					err := helperWhere.ExtractGreaterThan(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "hgt" {
					err := helperHaving.ExtractGreaterThan(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "gte" {
					err := helperWhere.ExtractGreaterThanOrEqual(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "hgte" {
					err := helperHaving.ExtractGreaterThanOrEqual(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "lt" {
					err := helperWhere.ExtractLowerThan(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "hlt" {
					err := helperHaving.ExtractLowerThan(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "lte" {
					err := helperWhere.ExtractLowerThanOrEqual(q, key, _rv)

					if err != nil {
						return err
					}
				} else if rk == "hlte" {
					err := helperHaving.ExtractLowerThanOrEqual(q, key, _rv)

					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (q *QSSQL) GetSQL() (string, error) {
	type Clause struct {
		Where   string
		OrderBy string
		Having  string
		Offset  string
		Limit   string
	}

	clause := &Clause{}
	clause.Where = q.GetWhereClause()
	clause.OrderBy = q.GetOrderByClause()
	clause.Having = q.GetHavingClause()
	clause.Offset = q.GetOffsetClause()
	clause.Limit = q.GetLimitClause()

	t, err := template.New("sql").Parse(q.template)

	if err != nil {
		return "", err
	}

	var b bytes.Buffer

	err = t.Execute(&b, clause)

	if err != nil {
		return "", nil
	}

	return b.String(), nil
}

func (q *QSSQL) AddWhereBindValue(param interface{}) {
	q.whereBindValues = append(q.whereBindValues, param)
}

func (q *QSSQL) GetWhereBindValues() []interface{} {
	return q.whereBindValues
}

func (q *QSSQL) AddHavingAlias(key string, value string) {
	q.havingAliases[key] = value
}

func (q *QSSQL) AddHavingAliases(aliases map[string]string) {
	for key, value := range aliases {
		q.havingAliases[key] = value
	}
}

func (q *QSSQL) GetHavingAliases() map[string]string {
	return q.havingAliases
}

func (q *QSSQL) AddHavingValue(value string) {
	q.havingValues = append(q.havingValues, value)
}

func (q *QSSQL) AddHavingValues(values []string) {
	for _, v := range values {
		q.havingValues = append(q.havingValues, v)
	}
}

func (q *QSSQL) GetHavingValues() []string {
	return q.havingValues
}

func (q *QSSQL) AddHavingBindValue(param interface{}) {
	q.havingBindValues = append(q.havingBindValues, param)
}

func (q *QSSQL) GetHavingBindValues() []interface{} {
	return q.havingBindValues
}

func (q *QSSQL) GetHavingClause() string {
	if len(q.havingValues) == 0 {
		return ""
	}

	return HAVING_KEYWORD + " " + strings.Join(q.havingValues[:], " AND ")
}

func (q *QSSQL) GetBindValues() []interface{} {
	var bindValues []interface{}

	bindValues = append(bindValues, q.whereBindValues...)
	bindValues = append(bindValues, q.havingBindValues...)

	return bindValues
}

func (q *QSSQL) GetLimitClause() string {
	if q.paramLimitValue > 0 {
		return fmt.Sprintf("%s %d", LIMIT_KEYWORD, q.paramLimitValue)
	}

	return ""
}

func (q *QSSQL) GetOffsetClause() string {
	if q.paramLimitValue >= 0 {
		return fmt.Sprintf("%s %d", OFFSET_KEYWORD, q.paramOffsetValue)
	}

	return ""
}

func (q *QSSQL) GetOrderByClause() string {
	if len(q.paramOrderByValues) == 0 {
		return ""
	}

	var orderByValues []string

	for _, v := range q.paramOrderByValues {
		if strings.HasPrefix(v, "-") {
			orderByValues = append(orderByValues, v[1:]+" DESC")
			continue
		}

		orderByValues = append(orderByValues, v)
	}

	orderByClause := ORDER_BY_KEYWORD + " "

	return orderByClause + strings.Join(orderByValues[:], ", ")
}
