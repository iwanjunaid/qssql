package offset

import (
	"errors"
	"fmt"

	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func Extract(q iface.GenericQSSQL, key string, value interface{}) error {
	_v, ok := value.(float64)

	if !ok {
		errorMessage := fmt.Sprintf("Param '%s' value should be a number", q.GetParamOffsetName())

		return errors.New(errorMessage)
	}

	if _v < 0 {
		errorMessage := fmt.Sprintf("Param '%s' value can't be a negative number", q.GetParamOffsetName())

		return errors.New(errorMessage)
	}

	q.SetParamOffsetValue(uint64(_v))

	return nil
}
