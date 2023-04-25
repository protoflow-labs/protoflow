package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResource(t *testing.T) {
	type TestStruct struct{}
	assert := assert.New(t)
	testStruct := &TestStruct{}

	testCases := []struct {
		name           string
		resources      map[string]any
		expectedResult any
		expectedError  error
	}{
		{
			name:           "Empty resources map",
			resources:      make(map[string]any),
			expectedResult: (*TestStruct)(nil),
			expectedError:  ErrResourceNotFound,
		},
		{
			name: "No matching resource",
			resources: map[string]any{
				"a": 1,
				"b": "hello",
			},
			expectedResult: (*TestStruct)(nil),
			expectedError:  ErrResourceNotFound,
		},
		{
			name: "Matching resource",
			resources: map[string]any{
				"a": "world",
				"b": 2,
				"c": testStruct,
			},
			expectedResult: testStruct,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getResource[TestStruct](tc.resources)

			assert.Equalf(tc.expectedError, err, "Test case '%s' failed. got: '%v', want: '%v'", tc.name, err, tc.expectedError)
			assert.Equalf(tc.expectedResult, result, "Test case '%s' failed. got: '%v', want: '%v'", tc.name, result, tc.expectedResult)
		})
	}

}
