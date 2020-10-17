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
		formatErrs, absenceErrs := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absenceErrs))
		if len(formatErrs) > 0 {
			assert.Equal(t, tc.wantIsExistErr, formatErrs[0].(Error).IsNotExist())
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
		formatErrs, absenceErrs := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absenceErrs))
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
				"age": 18,
			}},
			dataField:       "age",
			wantAbsence:     0,
			wantFormatError: 0,
		},
	}
	for _, tc := range cases {
		Check(tc.dataReq).Params(tc.dataField).IsInt32()
		formatErrs, absenceErrs := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absenceErrs))
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
		formatErrs, absenceErrs := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absenceErrs))
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
		formatErrs, absenceErrs := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absenceErrs))
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
		formatErrs, absenceErrs := ValidateResult(tc.dataReq)
		assert.Equal(t, tc.wantFormatError, len(formatErrs))
		assert.Equal(t, tc.wantAbsence, len(absenceErrs))
	}
}
