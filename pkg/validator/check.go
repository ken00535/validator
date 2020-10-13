package validator

import (
	"errors"
	"strconv"
)

// CheckType is type to sanitize
type CheckType struct {
	validatorBase
}

// Check return a check type to following operations
func Check(payload Payload) *CheckType {
	cache := payload.GetCache()
	if cache == nil {
		cache = make(map[string]interface{})
		payload.SetCache(cache)
	}
	if cache[contextKey] == nil {
		cache[contextKey] = make(map[string]interface{})
	}
	var errorList []error
	if cache[contextKey].(map[string]interface{})[errorsKey] != nil {
		errorList = cache[contextKey].(map[string]interface{})[errorsKey].([]error)
	}
	if errorList == nil {
		errorList = []error{}
		cache[contextKey].(map[string]interface{})[errorsKey] = errorList
	}
	var absence []string
	if cache[contextKey].(map[string]interface{})[abcenseKey] != nil {
		absence = cache[contextKey].(map[string]interface{})[abcenseKey].([]string)
	}
	if absence == nil {
		absence = []string{}
		cache[contextKey].(map[string]interface{})[abcenseKey] = absence
	}
	ret := &CheckType{}
	ret.content = payload
	return ret
}

// Params tag the param that will be sanitized
func (v *CheckType) Params(param string) *CheckType {
	v.param = param
	return v
}

// IsExist check param is exist or not
func (v *CheckType) IsExist() *CheckType {
	v.handleAbsence()
	return v
}

// IsInt check param is int or not
func (v *CheckType) IsInt() *CheckType {
	return v.isType(intType)
}

// IsUint32 check param is uint32 or not
func (v *CheckType) IsUint32() *CheckType {
	return v.isType(uint32Type)
}

// IsFloat check param is float64 or not
func (v *CheckType) IsFloat() *CheckType {
	return v.isType(float64Type)
}

// IsBool check param is bool or not
func (v *CheckType) IsBool() *CheckType {
	return v.isType(boolType)
}

func (v *CheckType) isType(dataType string) *CheckType {
	val, exist := v.handleAbsence()
	if exist {
		var err error
		switch dataType {
		case intType:
			_, err = strconv.Atoi(val)
		case uint32Type:
			_, err = strconv.ParseUint(val, 10, 32)
		case float64Type:
			_, err = strconv.ParseFloat(val, 64)
		case boolType:
			_, err = strconv.ParseBool(val)
		}
		v.handleErrors(err)
	}
	return v
}

func (v *CheckType) handleAbsence() (string, bool) {
	return handleAbsence(v)
}

func (v *CheckType) getAbsenceError() error {
	return errors.New(v.param + " don't exist!")
}