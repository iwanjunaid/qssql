package orderby

import (
	"errors"
	"fmt"

	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func Extract(q iface.GenericQSSQL, key string, value interface{}) error {
	switch value.(type) {
	case []interface{}:
		assertedValue := value.([]interface{})

		for _, val := range assertedValue {
			v, ok := val.(string)

			if !ok {
				errorMessage := fmt.Sprintf("Param '%s' value consist of non-string value", q.GetParamOrderByName())

				return errors.New(errorMessage)
			}

			q.AddParamOrderByValue(v)
		}
	default:
		errorMessage := fmt.Sprintf("Param '%s' should be in form of an array", q.GetParamOrderByName())

		return errors.New(errorMessage)
	}

	return nil
}
