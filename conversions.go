package forge

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

func asBoolean(value interface{}) (bool, error) {
	switch val := value.(type) {
	case bool:
		return val, nil
	case float64:
		return val != 0, nil
	case int64:
		return val != 0, nil
	case nil:
		return false, nil
	case string:
		return val != "", nil
	}

	msg := fmt.Sprintf("Could not convert value %s to type BOOLEAN", value)
	return false, errors.New(msg)
}

func asFloat(value interface{}) (float64, error) {
	switch val := value.(type) {
	case bool:
		if val {
			return float64(1), nil
		} else {
			return float64(0), nil
		}
	case float64:
		return val, nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	}

	msg := fmt.Sprintf("Could not convert value %s to type FLOAT", value)
	return 0, errors.New(msg)
}

func asInteger(value interface{}) (int64, error) {
	switch val := value.(type) {
	case bool:
		if val {
			return int64(1), nil
		} else {
			return int64(0), nil
		}
	case float64:
		return int64(math.Trunc(val)), nil
	case int64:
		return val, nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	}

	msg := fmt.Sprintf("Could not convert value %s to type INTEGER", value)
	return 0, errors.New(msg)
}

func asString(value interface{}) (string, error) {
	switch val := value.(type) {
	case bool:
		if val {
			return "True", nil
		} else {
			return "False", nil
		}
	case float64:
		return strconv.FormatFloat(val, 10, -1, 64), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case nil:
		return "Null", nil
	case string:
		return val, nil
	}

	msg := fmt.Sprintf("Could not convert value %s to type STRING", value)
	return "", errors.New(msg)
}
