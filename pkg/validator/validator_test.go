package validator

import (
	"testing"

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

type person struct {
	Name        string                 `vld:"name"`
	Gender      string                 `vld:"gender"`
	Age         int                    `vld:"age"`
	Score       int                    `vld:"score"`
	Weight      float64                `vld:"w"`
	IsAlive     bool                   `vld:"alive"`
	Description string                 `vld:"desc"`
	Hand        map[string]interface{} `vld:"hand"`
	Leg         leg                    `vld:"leg"`
	Parent      []string               `vld:"parent"`
	Houses      []house                `vld:"houses"`
}

func TestStruct(t *testing.T) {
	payload := &message{}
	expect := &person{Name: "ken"}
	Assign(payload).Struct(expect)
	actual := payload.GetCache()[contextKey].(map[string]interface{})["struct"]
	assert.Equal(t, expect, actual)
}

func TestSanitizeInt(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["age"] = "18"
	expect := person{Age: 18}
	actual := person{}
	Assign(payload).Struct(&actual)
	Sanitize(payload).Params("age").ToInt()
	assert.Equal(t, expect, actual)
}

func TestSanitizeBool(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["alive"] = "true"
	expect := person{IsAlive: true}
	actual := person{}
	Assign(payload).Struct(&actual)
	Sanitize(payload).Params("alive").ToBool()
	assert.Equal(t, expect, actual)
}

func TestSanitizeFloat(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["w"] = "64.5"
	expect := person{Weight: 64.5}
	actual := person{}
	Assign(payload).Struct(&actual)
	Sanitize(payload).Params("w").ToFloat64()
	assert.Equal(t, expect, actual)
}

func TestSanitizeString(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["desc"] = `"hello"`
	expect := person{Description: "hello"}
	actual := person{}
	Assign(payload).Struct(&actual)
	Sanitize(payload).Params("desc").ToString()
	assert.Equal(t, expect, actual)
}

func TestSanitizeObject(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      person
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"hand": `{"finger": 5}`,
			}},
			dataField: "hand",
			want:      person{Hand: map[string]interface{}{"finger": float64(5)}},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"hand": `{"finger": true}`,
			}},
			dataField: "hand",
			want:      person{Hand: map[string]interface{}{"finger": true}},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"hand": `{"finger": "this is my finger"}`,
			}},
			dataField: "hand",
			want:      person{Hand: map[string]interface{}{"finger": "this is my finger"}},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"hand": `{"finger": [1,2]}`,
			}},
			dataField: "hand",
			want:      person{Hand: map[string]interface{}{"finger": []interface{}{float64(1), float64(2)}}},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"hand": `{"finger": null}`,
			}},
			dataField: "hand",
			want:      person{Hand: map[string]interface{}{"finger": nil}},
		},
	}
	for _, tc := range cases {
		actual := person{}
		Assign(tc.dataReq).Struct(&actual)
		Sanitize(tc.dataReq).Params(tc.dataField).ToObject()
		assert.Equal(t, tc.want, actual)
	}
}

func TestSanitizeStruct(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      person
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"leg": `{"number": 2}`,
			}},
			dataField: "leg",
			want:      person{Leg: leg{Number: 2}},
		},
	}
	for _, tc := range cases {
		actual := person{}
		Assign(tc.dataReq).Struct(&actual)
		Sanitize(tc.dataReq).Params(tc.dataField).ToStruct()
		assert.Equal(t, tc.want, actual)
	}
}

func TestSanitizeSlice(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      person
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"parent": `["Mary", "Peter"]`,
			}},
			dataField: "parent",
			want:      person{Parent: []string{"Mary", "Peter"}},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"houses": `[{"size": 10, "win": 2}, {"size": 50, "win": 10}]`,
			}},
			dataField: "houses",
			want:      person{Houses: []house{{Size: 10, Window: 2}, {Size: 50, Window: 10}}},
		},
	}
	for _, tc := range cases {
		actual := person{}
		Assign(tc.dataReq).Struct(&actual)
		Sanitize(tc.dataReq).Params(tc.dataField).ToStruct()
		assert.Equal(t, tc.want, actual)
	}
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
				"age": "18",
			}},
			dataField:       "age",
			wantAbsence:     0,
			wantFormatError: 0,
		},
	}
	for _, tc := range cases {
		actual := person{}
		Assign(tc.dataReq).Struct(&actual)
		Sanitize(tc.dataReq).Params(tc.dataField).ToInt()
		formatErrs, absenceErrs := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absenceErrs))
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
		formatErrs, absenceErrs := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absenceErrs))
		Check(tc.dataReq).Params(tc.dataField2).IsExist()
		formatErrs, absenceErrs = ValidateResult(tc.dataReq)
		assert.Equal(t, 0, len(formatErrs))
		assert.Equal(t, 0, len(absenceErrs))
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
		actual := person{}
		Assign(tc.dataReq).Struct(&actual)
		Sanitize(tc.dataReq).Optional().Params(tc.dataField).ToInt()
		formatErrs, absenceErrs := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absenceErrs))
		if len(formatErrs) > 0 {
			assert.Equal(t, tc.wantIsExistErr, formatErrs[0].(Error).IsNotExist())
		}
	}
}

func TestTrim(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["age"] = "!!18"
	expect := person{Age: 18}
	actual := person{}
	Assign(payload).Struct(&actual)
	Sanitize(payload).Params("age").Trim("!").ToInt()
	assert.Equal(t, expect, actual)
}

func TestAnalyze(t *testing.T) {
	person := person{Name: "ken"}
	expect := []string{"Name"}
	tags := []string{"name"}
	actual := Analyze(person).Fields(tags)
	assert.Equal(t, expect, actual)
}
