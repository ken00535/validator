package validator

import (
	"reflect"
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

// IsInt32 check param is int32 or not
func (v *CheckType) IsInt32() *CheckType {
	return v.isType(int32Type)
}

// IsInt64 check param is int64 or not
func (v *CheckType) IsInt64() *CheckType {
	return v.isType(int64Type)
}

// IsString check param is string or not
func (v *CheckType) IsString() *CheckType {
	return v.isType(stringType)
}

// IsUint64 check param is uint64 or not
func (v *CheckType) IsUint64() *CheckType {
	return v.isType(uint64Type)
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

// IsBytes check param is bytes or not
func (v *CheckType) IsBytes() *CheckType {
	return v.isType(bytesType)
}

func (v *CheckType) isType(dataType int) *CheckType {
	val, exist := v.handleAbsence()
	if exist {
		var err error
		switch dataType {
		case intType:
			if reflect.TypeOf(val).Kind() != reflect.Int {
				err = newWrongTypeError("type is not int")
			}
		case int32Type:
			if reflect.TypeOf(val).Kind() != reflect.Int32 {
				err = newWrongTypeError("type is not int32")
			}
		case int64Type:
			if reflect.TypeOf(val).Kind() != reflect.Int64 {
				err = newWrongTypeError("type is not int64")
			}
		case uint32Type:
			if reflect.TypeOf(val).Kind() != reflect.Uint32 {
				err = newWrongTypeError("type is not uint32")
			}
		case uint64Type:
			if reflect.TypeOf(val).Kind() != reflect.Uint64 {
				err = newWrongTypeError("type is not uint64")
			}
		case float64Type:
			if reflect.TypeOf(val).Kind() != reflect.Float64 {
				err = newWrongTypeError("type is not float64")
			}
		case boolType:
			if reflect.TypeOf(val).Kind() != reflect.Bool {
				err = newWrongTypeError("type is not bool")
			}
		case stringType:
			if reflect.TypeOf(val).Kind() != reflect.String {
				err = newWrongTypeError("type is not string")
			}
		case bytesType:
			if reflect.TypeOf(val).Kind() != reflect.Slice ||
				reflect.TypeOf(val).Elem().Kind() != reflect.Uint8 {
				err = newWrongTypeError("type is not bytes")
			}
		}
		v.handleErrors(err)
	}
	return v
}

func (v *CheckType) handleAbsence() (interface{}, bool) {
	return handleAbsence(v)
}

func (v *CheckType) getAbsenceError() error {
	return newNotExistError(v.param + " don't exist!")
}
