package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeIntNew(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["age"] = "18"
	expect := person{Age: 18}
	actual := person{}
	SanitizeNew(payload).Params("age").ToIntNew(&actual)
	assert.Equal(t, expect, actual)
}

func TestSanitizeBoolNew(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["alive"] = "true"
	expect := person{IsAlive: true}
	actual := person{}
	SanitizeNew(payload).Params("alive").ToBoolNew(&actual)
	assert.Equal(t, expect, actual)
}

func TestSanitizeFloatNew(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["w"] = "64.5"
	expect := person{Weight: 64.5}
	actual := person{}
	SanitizeNew(payload).Params("w").ToFloat64New(&actual)
	assert.Equal(t, expect, actual)
}

func TestSanitizeStringNew(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      person
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"desc": `"hello"`,
			}},
			dataField: "desc",
			want:      person{Description: "hello"},
		},
	}
	for _, tc := range cases {
		actual := person{}
		SanitizeNew(tc.dataReq).Params(tc.dataField).ToStringNew(&actual)
		assert.Equal(t, tc.want, actual)
	}
}

func TestSanitizeObjectNew(t *testing.T) {
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
		SanitizeNew(tc.dataReq).Params(tc.dataField).ToObjectNew(&actual)
		assert.Equal(t, tc.want, actual)
	}
}

// func TestSanitizeStruct(t *testing.T) {
// 	type testCase struct {
// 		dataReq   *message
// 		dataField string
// 		want      person
// 	}
// 	cases := []testCase{
// 		{
// 			dataReq: &message{msg: map[string]interface{}{
// 				"leg": `{"number": 2}`,
// 			}},
// 			dataField: "leg",
// 			want:      person{Leg: leg{Number: 2}},
// 		},
// 	}
// 	for _, tc := range cases {
// 		actual := person{}
// 		Assign(tc.dataReq).Struct(&actual)
// 		Sanitize(tc.dataReq).Params(tc.dataField).ToStruct()
// 		assert.Equal(t, tc.want, actual)
// 	}
// }

// func TestSanitizeSlice(t *testing.T) {
// 	type testCase struct {
// 		dataReq   *message
// 		dataField string
// 		want      person
// 	}
// 	cases := []testCase{
// 		{
// 			dataReq: &message{msg: map[string]interface{}{
// 				"parent": `["Mary", "Peter"]`,
// 			}},
// 			dataField: "parent",
// 			want:      person{Parent: []string{"Mary", "Peter"}},
// 		},
// 		{
// 			dataReq: &message{msg: map[string]interface{}{
// 				"houses": `[{"size": 10, "win": 2}, {"size": 50, "win": 10}]`,
// 			}},
// 			dataField: "houses",
// 			want:      person{Houses: []house{{Size: 10, Window: 2}, {Size: 50, Window: 10}}},
// 		},
// 	}
// 	for _, tc := range cases {
// 		actual := person{}
// 		Assign(tc.dataReq).Struct(&actual)
// 		Sanitize(tc.dataReq).Params(tc.dataField).ToStruct()
// 		assert.Equal(t, tc.want, actual)
// 	}
// }
