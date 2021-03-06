package validator

import (
	"reflect"
)

const (
	contextKey = "github.com/govalidator/validator"
	abcenseKey = "abcense"
	errorsKey  = "errors"
	structKey  = "struct"
)

const (
	intType = iota
	int32Type
	int64Type
	uint32Type
	uint64Type
	float64Type
	boolType
	objectType
	stringType
	bytesType
	ipType
	timeType
	localTimeType
)

// Payload is payload of message, it will store some info of validator, so you have to
// create a map for it.
type Payload interface {
	GetCache() map[string]interface{}
	SetCache(map[string]interface{})
	GetParam(string) (val interface{}, exist bool)
}

// ValidateResult validate result of sanitize
func ValidateResult(payload Payload) (formatError []error, absence []string) {
	cache := payload.GetCache()
	errorList := cache[contextKey].(map[string]interface{})[errorsKey].([]error)
	absenceList := cache[contextKey].(map[string]interface{})[abcenseKey].([]string)
	cache[contextKey] = make(map[string]interface{})
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
