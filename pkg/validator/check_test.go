package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckExist(t *testing.T) {
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
			wantFormatError: 1,
			wantIsExistErr:  true,
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
		Check(tc.dataReq).Params(tc.dataField).IsExist()
		formatErrs, absence := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
		if len(formatErrs) > 0 {
			// assert.Equal(t, tc.wantIsExistErr, formatErrs[0].(Error).IsNotExist())
		}
	}
}

func TestCheckInt(t *testing.T) {
	type testCase struct {
		dataReq         *message
		dataField       string
		wantAbsence     int
		wantFormatError int
	}
	cases := []testCase{
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
				"age": 18,
			}},
			dataField:       "age",
			wantAbsence:     0,
			wantFormatError: 0,
		},
	}
	for _, tc := range cases {
		Check(tc.dataReq).Params(tc.dataField).IsInt()
		formatErrs, absence := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
	}
}

func TestCheckInt32(t *testing.T) {
	type testCase struct {
		dataReq         *message
		dataField       string
		wantAbsence     int
		wantFormatError int
	}
	cases := []testCase{
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
				"age": int32(18),
			}},
			dataField:       "age",
			wantAbsence:     0,
			wantFormatError: 0,
		},
	}
	for _, tc := range cases {
		Check(tc.dataReq).Params(tc.dataField).IsInt32()
		formatErrs, absence := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
	}
}

func TestCheckBytes(t *testing.T) {
	type testCase struct {
		dataReq         *message
		dataField       string
		wantAbsence     int
		wantFormatError int
	}
	cases := []testCase{
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
				"age": []byte("18A"),
			}},
			dataField:       "age",
			wantAbsence:     0,
			wantFormatError: 0,
		},
	}
	for _, tc := range cases {
		Check(tc.dataReq).Params(tc.dataField).IsBytes()
		formatErrs, absence := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
	}
}

func TestCheckBool(t *testing.T) {
	type testCase struct {
		dataReq         *message
		dataField       string
		wantAbsence     int
		wantFormatError int
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"alive": "true",
			}},
			dataField:       "alive",
			wantAbsence:     0,
			wantFormatError: 1,
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"alive": true,
			}},
			dataField:       "alive",
			wantAbsence:     0,
			wantFormatError: 0,
		},
	}
	for _, tc := range cases {
		Check(tc.dataReq).Params(tc.dataField).IsBool()
		formatErrs, absence := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
	}
}

func TestCheckFloat(t *testing.T) {
	type testCase struct {
		dataReq         *message
		dataField       string
		wantAbsence     int
		wantFormatError int
	}
	cases := []testCase{
		{
			dataReq: &message{msg: map[string]interface{}{
				"w": "67.1",
			}},
			dataField:       "w",
			wantAbsence:     0,
			wantFormatError: 1,
		},
		{
			dataReq: &message{msg: map[string]interface{}{
				"w": 67.1,
			}},
			dataField:       "w",
			wantAbsence:     0,
			wantFormatError: 0,
		},
	}
	for _, tc := range cases {
		Check(tc.dataReq).Params(tc.dataField).IsFloat()
		formatErrs, absence := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absence))
	}
}
