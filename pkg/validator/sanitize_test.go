package validator

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeInt(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      testStruct
	}
	hp := 180
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18",
			}},
			dataField: "age",
			want:      testStruct{Age: 18},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"age": "18",
			}},
			dataField: "hp",
			want:      testStruct{HP: nil},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"hp": "180",
			}},
			dataField: "hp",
			want:      testStruct{HP: &hp},
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField).ToInt(&actual)
		assert.Equal(t, tc.want, actual)
	}
}

func TestSanitizeBool(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["alive"] = "true"
	expect := testStruct{IsAlive: true}
	actual := testStruct{}
	Sanitize(payload).Params("alive").ToBool(&actual)
	assert.Equal(t, expect, actual)
}

func TestSanitizeFloat(t *testing.T) {
	payload := &message{msg: map[string]interface{}{}}
	payload.msg["w"] = "64.5"
	expect := testStruct{Weight: 64.5}
	actual := testStruct{}
	Sanitize(payload).Params("w").ToFloat64(&actual)
	assert.Equal(t, expect, actual)
}

func TestSanitizeString(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      testStruct
	}
	strPtr := "hello"
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"desc": `"hello"`,
			}},
			dataField: "desc",
			want:      testStruct{Description: "hello"},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"strptr": `"hello"`,
			}},
			dataField: "strptr",
			want:      testStruct{StrPtr: &strPtr},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"null": `"hello"`,
			}},
			dataField: "strptr",
			want:      testStruct{StrPtr: nil},
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField).ToString(&actual)
		assert.Equal(t, tc.want, actual)
	}
}

func TestSanitizeObject(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      testStruct
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"hand": `{"finger": 5}`,
			}},
			dataField: "hand",
			want:      testStruct{Hand: map[string]interface{}{"finger": float64(5)}},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"hand": `{"finger": true}`,
			}},
			dataField: "hand",
			want:      testStruct{Hand: map[string]interface{}{"finger": true}},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"hand": `{"finger": "this is my finger"}`,
			}},
			dataField: "hand",
			want:      testStruct{Hand: map[string]interface{}{"finger": "this is my finger"}},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"hand": `{"finger": [1,2]}`,
			}},
			dataField: "hand",
			want:      testStruct{Hand: map[string]interface{}{"finger": []interface{}{float64(1), float64(2)}}},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"hand": `{"finger": null}`,
			}},
			dataField: "hand",
			want:      testStruct{Hand: map[string]interface{}{"finger": nil}},
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField).ToObject(&actual)
		assert.Equal(t, tc.want, actual)
	}
}

func TestSanitizeStruct(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      testStruct
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"leg": `{"number": 2}`,
			}},
			dataField: "leg",
			want:      testStruct{Leg: leg{Number: 2}},
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField).ToStruct(&actual)
		assert.Equal(t, tc.want, actual)
	}
}

func TestSanitizeSlice(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      testStruct
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"parent": `["Mary", "Peter"]`,
			}},
			dataField: "parent",
			want:      testStruct{Parent: []string{"Mary", "Peter"}},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"houses": `[{"size": 10, "win": 2}, {"size": 50, "win": 10}]`,
			}},
			dataField: "houses",
			want:      testStruct{Houses: []house{{Size: 10, Window: 2}, {Size: 50, Window: 10}}},
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField).ToStruct(&actual)
		assert.Equal(t, tc.want, actual)
	}
}

func TestSanitizeIP(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      testStruct
	}
	ipPtr := net.IPv4(127, 0, 0, 1)
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"ip": "127.0.0.1",
			}},
			dataField: "ip",
			want:      testStruct{IP: net.IPv4(127, 0, 0, 1)},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"ipPtr": "127.0.0.1",
			}},
			dataField: "ipPtr",
			want:      testStruct{IPPtr: &ipPtr},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"null": "127.0.0.1",
			}},
			dataField: "ipPtr",
			want:      testStruct{IPPtr: nil},
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField).ToIP(&actual)
		assert.Equal(t, tc.want, actual)
	}
}

func TestSanitizeLocalTime(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      testStruct
	}
	timeVal, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-11-06 16:19:23", time.Local)
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"startTime": "2020-11-06 16:19:23",
			}},
			dataField: "startTime",
			want:      testStruct{StartTime: timeVal},
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField).TimeFormat("2006-01-02 15:04:05").ToLocalTime(&actual)
		assert.Equal(t, tc.want, actual)
	}
}

func TestSanitizeTime(t *testing.T) {
	type testCase struct {
		dataReq   *message
		dataField string
		want      testStruct
	}
	timeVal, _ := time.Parse("2006-01-02 15:04:05", "2020-11-06 16:19:23")
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"startTime": "2020-11-06 16:19:23",
			}},
			dataField: "startTime",
			want:      testStruct{StartTime: timeVal},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"endTime": "2020-11-06 16:19:23",
			}},
			dataField: "endTime",
			want:      testStruct{EndTime: &timeVal},
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"startTime": "2020-11-06 16:19:23",
			}},
			dataField: "endTime",
			want:      testStruct{EndTime: nil},
		},
	}
	for _, tc := range cases {
		actual := testStruct{}
		Sanitize(tc.dataReq).Params(tc.dataField).TimeFormat("2006-01-02 15:04:05").ToTime(&actual)
		assert.Equal(t, tc.want, actual)
	}
}
