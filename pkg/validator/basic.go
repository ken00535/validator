package validator

type validatorInterface interface {
	getAbsenceError() error
	getPayload() Payload
	getOptional() bool
	getParam() string
}

type validatorBase struct {
	content  Payload
	param    string
	optional bool
}

func (v *validatorBase) getPayload() Payload {
	return v.content
}

func (v *validatorBase) getOptional() bool {
	return v.optional
}

func (v *validatorBase) getParam() string {
	return v.param
}

func (v *validatorBase) handleErrors(err error) {
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

func handleAbsence(v validatorInterface) (interface{}, bool) {
	cache := v.getPayload().GetCache()
	errorList := cache[contextKey].(map[string]interface{})[errorsKey].([]error)
	absenceList := cache[contextKey].(map[string]interface{})[abcenseKey].([]string)
	val, exist := v.getPayload().GetParam(v.getParam())
	if !exist {
		if !v.getOptional() {
			errorList = append(errorList, v.getAbsenceError())
			cache[contextKey].(map[string]interface{})[errorsKey] = errorList
		}
		absenceList = append(absenceList, v.getParam())
		cache[contextKey].(map[string]interface{})[abcenseKey] = absenceList
	}
	return val, exist
}
