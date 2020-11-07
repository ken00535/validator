package validator

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type message struct {
	msg   map[string]interface{}
	cache map[string]interface{}
}

func (m *message) GetCache() map[string]interface{} {
	if m.cache == nil {
		m.cache = make(map[string]interface{})
	}
	return m.cache
}

func (m *message) SetCache(input map[string]interface{}) {
	m.cache = input
}

func (m *message) GetParam(field string) (val interface{}, exist bool) {
	v, ok := m.msg[field]
	return v, ok
}

type leg struct {
	Number int `json:"number"`
}

type house struct {
	Size   int `json:"size"`
	Window int `json:"win"`
}

type testStruct struct {
	Name        string                 `vld:"name"`
	Gender      string                 `vld:"gender"`
	Age         int                    `vld:"age"`
	HP          *int                   `vld:"hp"`
	Score       int                    `vld:"score"`
	Weight      float64                `vld:"w"`
	IsAlive     bool                   `vld:"alive"`
	Description string                 `vld:"desc"`
	Hand        map[string]interface{} `vld:"hand"`
	Leg         leg                    `vld:"leg"`
	Parent      []string               `vld:"parent"`
	Houses      []house                `vld:"houses"`
	IP          net.IP                 `vld:"ip"`
	StartTime   time.Time              `vld:"startTime"`
	EndTime     *time.Time             `vld:"endTime"`
}

func TestValidateResult(t *testing.T) {
	type testCase struct {
		dataReq         *message
		dataField       string
		wantAbsence     int
		wantFormatError int
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18",
			}},
			dataField:       "score",
			wantAbsence:     1,
			wantFormatError: 1,
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18A",
			}},
			dataField:       "age",
			wantAbsence:     0,
			wantFormatError: 1,
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18",
			}},
			dataField:       "age",
			wantAbsence:     0,
			wantFormatError: 0,
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField).ToInt(&actual)
		formatErrs, absence := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
	}
}

func TestValidateError(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		wantError error
		wantMsg   string
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18",
			}},
			dataField: "score",
			wantError: NotExistError{},
			wantMsg:   "score",
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18A",
			}},
			dataField: "age",
			wantError: WrongTypeError{},
			wantMsg:   "age",
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18",
			}},
			dataField: "age",
			wantError: nil,
			wantMsg:   "",
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField).ToInt(&actual)
		errs, _ := ValidateResult(tc.dataReq)
		if len(errs) > 0 {
			assert.True(t, reflect.TypeOf(errs[0]).AssignableTo(reflect.TypeOf(tc.wantError)))
			assert.Contains(t, errs[0].Error(), tc.wantMsg)
		}
	}
}

func TestValidateClear(t *testing.T) {
	type testCase struct {
		dataReq         *message
		dataField1      string
		dataField2      string
		wantAbsence     int
		wantFormatError int
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18",
			}},
			dataField1:      "score",
			dataField2:      "age",
			wantAbsence:     1,
			wantFormatError: 1,
		},
	}
	for _, tc := range cases {
		Check(tc.dataReq).Params(tc.dataField1).IsExist()
		formatErrs, absence := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
		Check(tc.dataReq).Params(tc.dataField2).IsExist()
		formatErrs, absence = ValidateResult(tc.dataReq)
		assert.Equal(t, 0, len(formatErrs))
		assert.Equal(t, 0, len(absence))
	}
}

func TestValidateErrorCount(t *testing.T) {
	type testCase struct {
		dataReq         *message
		dataField1      string
		dataField2      string
		wantAbsence     int
		wantFormatError int
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18",
			}},
			dataField1:      "score",
			dataField2:      "leg",
			wantAbsence:     2,
			wantFormatError: 2,
		},
	}
	for _, tc := range cases {
		Check(tc.dataReq).Params(tc.dataField1).IsExist()
		Check(tc.dataReq).Params(tc.dataField2).IsExist()
		formatErrs, absence := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
		person := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField1).ToInt(&person)
		Sanitize(tc.dataReq).Params(tc.dataField2).ToInt(&person)
		formatErrs, absence = ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
	}
}

func TestValidateOptional(t *testing.T) {
	type testCase struct {
		dataReq         *message
		dataField       string
		wantAbsence     int
		wantFormatError int
		wantIsExistErr  bool
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18",
			}},
			dataField:       "score",
			wantAbsence:     1,
			wantFormatError: 0,
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Optional().Params(tc.dataField).ToInt(&actual)
		formatErrs, absence := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
		if len(formatErrs) > 0 {
			// assert.Equal(t, tc.wantIsExistErr, formatErrs[0].(Error).IsNotExist())
		}
	}
}

func TestTrim(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["age"] = "!!18"
	expect := testStruct{Age: 18}
	actual := testStruct{}
	Sanitize(payload).Params("age").Trim("!").ToInt(&actual)
	assert.Equal(t, expect, actual)
}

func TestAnalyze(t *testing.T) {
	person := testStruct{Name: "ken"}
	expect := []string{"Name"}
	tags := []string{"name"}
	actual := Analyze(person).Fields(tags)
	assert.Equal(t, expect, actual)
}
