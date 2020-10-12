package validator

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

const (
	contextKey = "github.com/govalidator/validator"
	abcenseKey = "abcense"
	errorsKey  = "errors"
	structKey  = "struct"
)

const (
	intType     = "int"
	uint32Type  = "uint32"
	float64Type = "float64"
	boolType    = "bool"
	objectType  = "object"
)

// Payload is payload of message, it will store some info of validator, so you have to
// create a map for it.
type Payload interface {
	GetCache() map[string]interface{}
	SetCache(map[string]interface{})
	GetParam(string) (val string, exist bool)
}

// AssignType is type to validate
type AssignType struct {
	content Payload
}

// Assign a message to validate
func Assign(payload Payload) *AssignType {
	return &AssignType{
		content: payload,
	}
}

// Struct assign instance that store the message
func (v *AssignType) Struct(target interface{}) *AssignType {
	cache := v.content.GetCache()
	if cache == nil {
		cache = make(map[string]interface{})
		v.content.SetCache(cache)
	}
	cache[contextKey] = make(map[string]interface{})
	cache[contextKey].(map[string]interface{})[structKey] = target
	return v
}

// SanitizeType is type to sanitize
type SanitizeType struct {
	content  Payload
	param    string
	optional bool
	cutset   string
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
	return &SanitizeType{
		content: payload,
	}
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

func (v *SanitizeType) toValue(dataType string) *SanitizeType {
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
		case uint32Type:
			var uint32Instance uint64
			uint32Instance, err = strconv.ParseUint(val, 10, 32)
			valInstance = uint32(uint32Instance)
			field.Set(reflect.ValueOf(valInstance))
		case float64Type:
			valInstance, err = strconv.ParseFloat(val, 64)
			field.Set(reflect.ValueOf(valInstance))
		case boolType:
			valInstance, err = strconv.ParseBool(val)
			field.Set(reflect.ValueOf(valInstance))
		case objectType:
			varAddr := field.Addr().Interface()
			err = json.Unmarshal([]byte(val), varAddr)
		}
		v.handleErrors(err)
	}
	return v
}

func (v *SanitizeType) handleAbsence() (string, bool) {
	cache := v.content.GetCache()
	errorList := cache[contextKey].(map[string]interface{})[errorsKey].([]error)
	absenceList := cache[contextKey].(map[string]interface{})[abcenseKey].([]string)
	val, exist := v.content.GetParam(v.param)
	if !exist {
		if !v.optional {
			errorList = append(errorList, errors.New("lack of "+v.param))
			cache[contextKey].(map[string]interface{})[errorsKey] = errorList
		}
		absenceList = append(absenceList, v.param)
		cache[contextKey].(map[string]interface{})[abcenseKey] = absenceList
	}
	return val, exist
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

// ValidateResult validate result of sanitize
func ValidateResult(payload Payload) (formatError []error, absence []string) {
	cache := payload.GetCache()
	errorList := cache[contextKey].(map[string]interface{})[errorsKey].([]error)
	absenceList := cache[contextKey].(map[string]interface{})[abcenseKey].([]string)
	return errorList, absenceList
}

// AnalyzeType is type to validate
type AnalyzeType struct {
	content interface{}
}

// Analyze a struct
func Analyze(structTagged interface{}) *AnalyzeType {
	return &AnalyzeType{
		content: structTagged,
	}
}

// Fields get the tagged field name
func (v *AnalyzeType) Fields(tags []string) []string {
	tagsMap := make(map[string]bool)
	for _, tag := range tags {
		tagsMap[tag] = true
	}
	fieldNames := []string{}
	contentType := reflect.TypeOf(v.content)
	for i := 0; i < contentType.NumField(); i++ {
		if _, ok := tagsMap[contentType.Field(i).Tag.Get("vld")]; ok {
			fieldNames = append(fieldNames, contentType.Field(i).Name)
		}
	}
	return fieldNames
}

// CheckType is type to sanitize
type CheckType struct {
	content  Payload
	param    string
	optional bool
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
	return &CheckType{
		content: payload,
	}
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

func (v *CheckType) handleAbsence() (string, bool) {
	cache := v.content.GetCache()
	errorList := cache[contextKey].(map[string]interface{})[errorsKey].([]error)
	absenceList := cache[contextKey].(map[string]interface{})[abcenseKey].([]string)
	val, exist := v.content.GetParam(v.param)
	if !exist {
		if !v.optional {
			errorList = append(errorList, errors.New(v.param+" don't exist!"))
			cache[contextKey].(map[string]interface{})[errorsKey] = errorList
		}
		absenceList = append(absenceList, v.param)
		cache[contextKey].(map[string]interface{})[abcenseKey] = absenceList
	}
	return val, exist
}
