package limit

import (
	"errors"
	"fmt"
	"strconv"

	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func Extract(q iface.GenericQSSQL, key string, value interface{}) error {
	stringValue, ok := value.(string)

	if !ok {
		errorMessage := fmt.Sprintf("Can't assert to string for key %s", key)

		return errors.New(errorMessage)
	}

	floatParsed, err := strconv.ParseFloat(stringValue, 64)

	if err != nil {
		errorMessage := fmt.Sprintf("Param '%s' value should be a number", q.GetParamOffsetName())

		return errors.New(errorMessage)
	}

	if floatParsed < 1 {
		errorMessage := fmt.Sprintf("Param '%s' value should be a positive number", q.GetParamLimitName())

		return errors.New(errorMessage)
	}

	q.SetParamLimitValue(uint64(floatParsed))

	return nil
}
