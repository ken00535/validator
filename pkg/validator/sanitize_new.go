package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Sanitize return a sanitize type to following operations
func SanitizeNew(payload Payload) *SanitizeType {
	cache := payload.GetCache()
	if cache == nil {
		cache = make(map[string]interface{})
		payload.SetCache(cache)
	}
	cache[contextKey] = make(map[string]interface{})

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

// ToIntNew sanitize field to int
func (v *SanitizeType) ToIntNew(out interface{}) *SanitizeType {
	v.toValueNew(out, intType)
	return v
}

// ToBoolNew sanitize field to bool
func (v *SanitizeType) ToBoolNew(out interface{}) *SanitizeType {
	v.toValueNew(out, boolType)
	return v
}

func (v *SanitizeType) toValueNew(out interface{}, dataType int) *SanitizeType {
	val, exist := v.handleAbsence()
	if exist {
		if v.cutset != "" {
			val = strings.Trim(val, v.cutset)
		}
		var valInstance interface{}
		var field = v.getFieldNew(out)
		var err error
		switch dataType {
		case intType:
			valInstance, err = strconv.Atoi(val)
			field.Set(reflect.ValueOf(valInstance))
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not int", val))
			}
		case uint32Type:
			var uint32Instance uint64
			uint32Instance, err = strconv.ParseUint(val, 10, 32)
			valInstance = uint32(uint32Instance)
			field.Set(reflect.ValueOf(valInstance))
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not int32", val))
			}
		case float64Type:
			valInstance, err = strconv.ParseFloat(val, 64)
			field.Set(reflect.ValueOf(valInstance))
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not float", val))
			}
		case boolType:
			valInstance, err = strconv.ParseBool(val)
			field.Set(reflect.ValueOf(valInstance))
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not bool", val))
			}
		case objectType:
			varAddr := field.Addr().Interface()
			err = json.Unmarshal([]byte(val), varAddr)
			if err != nil {
				err = newWrongTypeError(fmt.Sprintf("message %v is not json or string", val))
			}
		}
		v.handleErrors(err)
	}
	return v
}

func (v *SanitizeType) getFieldNew(out interface{}) reflect.Value {
	// cache := v.content.GetCache()
	// target := cache[contextKey].(map[string]interface{})[structKey]
	targetValue := reflect.ValueOf(out).Elem()
	targetType := reflect.TypeOf(out).Elem()
	for i := 0; i < targetType.NumField(); i++ {
		if targetType.Field(i).Tag.Get("vld") == v.param {
			return targetValue.Field(i)
		}
	}
	return reflect.Value{}
}
