package validator

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// SanitizeType is type to sanitize
type SanitizeType struct {
	validatorBase
	cutset     string
	timeFormat string
}

// Sanitize return a sanitize type to following operations
func Sanitize(payload Payload) *SanitizeType {
	cache := payload.GetCache()
	if cache == nil {
		cache = make(map[string]interface{})
		payload.SetCache(cache)
	}
	if cache[contextKey] == nil {
		cache[contextKey] = make(map[string]interface{})
	}

	var errorList []error
	cache = payload.GetCache()
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

// ToInt sanitize field to int
func (v *SanitizeType) ToInt(out interface{}) *SanitizeType {
	v.toValue(out, intType)
	return v
}

// ToUint32 sanitize field to uint32
func (v *SanitizeType) ToUint32(out interface{}) *SanitizeType {
	v.toValue(out, uint32Type)
	return v
}

// ToBool sanitize field to bool
func (v *SanitizeType) ToBool(out interface{}) *SanitizeType {
	v.toValue(out, boolType)
	return v
}

// ToFloat64 sanitize field to float64
func (v *SanitizeType) ToFloat64(out interface{}) *SanitizeType {
	v.toValue(out, float64Type)
	return v
}

// ToObject sanitize field to object
func (v *SanitizeType) ToObject(out interface{}) *SanitizeType {
	v.toValue(out, objectType)
	return v
}

// ToString sanitize field to string
func (v *SanitizeType) ToString(out interface{}) *SanitizeType {
	v.toValue(out, objectType)
	return v
}

// ToStruct sanitize field to struct
func (v *SanitizeType) ToStruct(out interface{}) *SanitizeType {
	v.toValue(out, objectType)
	return v
}

// ToIP sanitize field to ip
func (v *SanitizeType) ToIP(out interface{}) *SanitizeType {
	v.toValue(out, ipType)
	return v
}

// ToTime sanitize field to time
func (v *SanitizeType) ToTime(out interface{}) *SanitizeType {
	v.toValue(out, timeType)
	return v
}

// ToLocalTime sanitize field to time
func (v *SanitizeType) ToLocalTime(out interface{}) *SanitizeType {
	v.toValue(out, localTimeType)
	return v
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

// TimeFormat provide time format info to sanitize
func (v *SanitizeType) TimeFormat(str string) *SanitizeType {
	v.timeFormat = str
	return v
}

func (v *SanitizeType) toValue(out interface{}, dataType int) *SanitizeType {
	val, exist := v.handleAbsence()
	if exist {
		if v.cutset != "" {
			val = strings.Trim(val, v.cutset)
		}
		var valInstance interface{}
		var field = v.getField(out)
		var err error
		switch dataType {
		case intType:
			valInstance, err = strconv.Atoi(val)
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not int", val))
			} else {
				setField(field, valInstance)
			}
		case uint32Type:
			var uint32Instance uint64
			uint32Instance, err = strconv.ParseUint(val, 10, 32)
			valInstance = uint32(uint32Instance)
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not int32", val))
			} else {
				setField(field, valInstance)
			}
		case float64Type:
			valInstance, err = strconv.ParseFloat(val, 64)
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not float", val))
			} else {
				setField(field, valInstance)
			}
		case boolType:
			valInstance, err = strconv.ParseBool(val)
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not bool", val))
			} else {
				setField(field, valInstance)
			}
		case objectType:
			varAddr := field.Addr().Interface()
			err = json.Unmarshal([]byte(val), varAddr)
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not json or string", val))
			}
		case ipType:
			valInstance = net.ParseIP(val)
			field.Set(reflect.ValueOf(valInstance))
			if valInstance == nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not ip", val))
			}
		case timeType:
			valInstance, err = time.Parse(v.timeFormat, val)
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not time. parse error: %v", val, err.Error()))
			} else {
				setField(field, valInstance)
			}
		case localTimeType:
			valInstance, err = time.ParseInLocation(v.timeFormat, val, time.Local)
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not time. parse error: %v", val, err.Error()))
			} else {
				setField(field, valInstance)
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

func (v *SanitizeType) getField(out interface{}) reflect.Value {
	targetValue := reflect.ValueOf(out).Elem()
	targetType := reflect.TypeOf(out).Elem()
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

func setField(field reflect.Value, val interface{}) {
	if field.Kind() == reflect.Ptr {
		switch val.(type) {
		case int:
			val := val.(int)
			field.Set(reflect.ValueOf(&val))
		case uint32:
			val := val.(uint32)
			field.Set(reflect.ValueOf(&val))
		case float64:
			val := val.(float64)
			field.Set(reflect.ValueOf(&val))
		case bool:
			val := val.(bool)
			field.Set(reflect.ValueOf(&val))
		case time.Time:
			val := val.(time.Time)
			field.Set(reflect.ValueOf(&val))
		}
	} else {
		field.Set(reflect.ValueOf(val))
	}
}
