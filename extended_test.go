package jsonata

import (
	"os"
	"testing"

	"github.com/goccy/go-json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChassis(t *testing.T) {
	tests := []struct {
		// the name of the test
		Name               string
		// the input file (must be JSON)
		InputFile          string
		// the jsonata that is applied to the file
		InputJsonataFile   string
		// the expected output file that is the result of the transformations
		ExpectedOutputFile string
	}{
		{
			Name:               "a simple test",
			InputFile:          "extendedTestFiles/case1/input.json",
			InputJsonataFile:   "extendedTestFiles/case1/input.jsonata",
			ExpectedOutputFile: "extendedTestFiles/case1/output.json",
		},
	}

	for index := range tests {
		tt := tests[index]

		t.Run(tt.Name, func(t *testing.T) {
			inputBytes, err := os.ReadFile(tt.InputFile)
			require.NoError(t, err)

			jsonataBytes, err := os.ReadFile(tt.InputJsonataFile)
			require.NoError(t, err)

			outputBytes, err := os.ReadFile(tt.ExpectedOutputFile)
			require.NoError(t, err)

			expr, err := Compile(string(jsonataBytes))
			require.NoError(t, err)

			var input, output interface{}

			err = json.Unmarshal(inputBytes, &input)
			require.NoError(t, err)

			result, err := expr.Eval(input)
			require.NoError(t, err)

			err = json.Unmarshal(outputBytes, &output)
			require.NoError(t, err)

			assert.Equal(t, result, output)
		})
	}
}
