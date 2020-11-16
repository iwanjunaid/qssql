package limit

import (
	"errors"
	"fmt"

	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func Extract(q iface.GenericQSSQL, key string, value interface{}) error {
	_v, ok := value.(float64)

	if !ok {
		errorMessage := fmt.Sprintf("Param '%s' value should be a number", q.GetParamLimitName())

		return errors.New(errorMessage)
	}

	if _v < 1 {
		errorMessage := fmt.Sprintf("Param '%s' value should be a positive number", q.GetParamLimitName())

		return errors.New(errorMessage)
	}

	q.SetParamLimitValue(uint64(_v))

	return nil
}
