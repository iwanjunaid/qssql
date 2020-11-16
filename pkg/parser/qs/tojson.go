package qs

import (
	"github.com/joncalhoun/qson"
)

func ToJSON(url string) (string, error) {
	res, err := qson.ToJSON(url)

	if err != nil {
		return "", err
	}

	return string(res), nil
}
