package validator

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// SanitizeType is type to sanitize
type SanitizeType struct {
	validatorBase
	cutset string
}

// Sanitize return a sanitize type to following operations
func Sanitize(payload Payload) *SanitizeType {
	var errorList []error
	cache := payload.GetCache()
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
	ret := &SanitizeType{}
	ret.content = payload
	return ret
}

// Params tag the param that will be sanitized
func (v *SanitizeType) Params(param string) *SanitizeType {
	v.param = param
	return v
}

// Optional tag the field is optinal
func (v *SanitizeType) Optional() *SanitizeType {
	v.optional = true
	return v
}

// Trim the unused part before assign
func (v *SanitizeType) Trim(str string) *SanitizeType {
	v.cutset = str
	return v
}

// ToInt sanitize field to int
func (v *SanitizeType) ToInt() *SanitizeType {
	v.toValue(intType)
	return v
}

// ToUint32 sanitize field to uint32
func (v *SanitizeType) ToUint32() *SanitizeType {
	v.toValue(uint32Type)
	return v
}

// ToFloat64 sanitize field to float64
func (v *SanitizeType) ToFloat64() *SanitizeType {
	v.toValue(float64Type)
	return v
}

// ToBool sanitize field to bool
func (v *SanitizeType) ToBool() *SanitizeType {
	v.toValue(boolType)
	return v
}

// ToObject sanitize field to object
func (v *SanitizeType) ToObject() *SanitizeType {
	v.toValue(objectType)
	return v
}

// ToStruct sanitize field to struct
func (v *SanitizeType) ToStruct() *SanitizeType {
	v.toValue(objectType)
	return v
}

// ToString sanitize field to string
func (v *SanitizeType) ToString() *SanitizeType {
	v.toValue(objectType)
	return v
}

func (v *SanitizeType) toValue(dataType int) *SanitizeType {
	val, exist := v.handleAbsence()
	if exist {
		if v.cutset != "" {
			val = strings.Trim(val, v.cutset)
		}
		var valInstance interface{}
		var field = v.getField()
		var err error
		switch dataType {
		case intType:
			valInstance, err = strconv.Atoi(val)
			field.Set(reflect.ValueOf(valInstance))
			if err != nil {
				err = newWrongTypeError("message is not int")
			}
		case uint32Type:
			var uint32Instance uint64
			uint32Instance, err = strconv.ParseUint(val, 10, 32)
			valInstance = uint32(uint32Instance)
			field.Set(reflect.ValueOf(valInstance))
			if err != nil {
				err = newWrongTypeError("message is not int32")
			}
		case float64Type:
			valInstance, err = strconv.ParseFloat(val, 64)
			field.Set(reflect.ValueOf(valInstance))
			if err != nil {
				err = newWrongTypeError("message is not float")
			}
		case boolType:
			valInstance, err = strconv.ParseBool(val)
			field.Set(reflect.ValueOf(valInstance))
			if err != nil {
				err = newWrongTypeError("message is not bool")
			}
		case objectType:
			varAddr := field.Addr().Interface()
			err = json.Unmarshal([]byte(val), varAddr)
			if err != nil {
				err = newWrongTypeError("message is not json or string")
			}
		}
		v.handleErrors(err)
	}
	return v
}

func (v *SanitizeType) handleAbsence() (string, bool) {
	val, exist := handleAbsence(v)
	valStr, ok := val.(string)
	return valStr, (exist && ok)
}

func (v *SanitizeType) handleErrors(err error) {
	cache := v.content.GetCache()
	var errorList []error
	if cache[contextKey].(map[string]interface{})[errorsKey] != nil {
		errorList = cache[contextKey].(map[string]interface{})[errorsKey].([]error)
	}
	if errorList == nil {
		errorList = []error{}
		cache[contextKey].(map[string]interface{})[errorsKey] = errorList
	}
	if err != nil {
		errorList = append(errorList, err)
		cache[contextKey].(map[string]interface{})[errorsKey] = errorList
	}
}

func (v *SanitizeType) getField() reflect.Value {
	cache := v.content.GetCache()
	target := cache[contextKey].(map[string]interface{})[structKey]
	targetValue := reflect.ValueOf(target).Elem()
	targetType := reflect.TypeOf(target).Elem()
	for i := 0; i < targetType.NumField(); i++ {
		if targetType.Field(i).Tag.Get("vld") == v.param {
			return targetValue.Field(i)
		}
	}
	return reflect.Value{}
}

func (v *SanitizeType) getAbsenceError() error {
	return newNotExistError(v.param + " don't exist!")
}
