package compute

import (
	"testing"
	"umemory/internal/compute"
)

type parserTestCase struct {
	name string
	arg string
	expectedCommand string
	expectedArgs []string
	expectedErrText string
}

func TestRequestParser(t *testing.T) {
	parser := compute.NewRequestParser()

	var testCases = []parserTestCase{
		{
			name: "get validate error",
			arg: "get asd qwe",
			expectedCommand: "",
			expectedArgs: nil,
			expectedErrText: "expected 1 argument, got 2",
		},
		{
			name: "set validate error",
			arg: "set asd",
			expectedCommand: "",
			expectedArgs: nil,
			expectedErrText: "expected 2 arguments, got 1",
		},
		{
			name: "delete validate error",
			arg: "delete asd qwe",
			expectedCommand: "",
			expectedArgs: nil,
			expectedErrText: "expected 1 argument, got 2",
		},
		{
			name: "unexpected command error",
			arg: "asdasd qwe",
			expectedCommand: "",
			expectedArgs: nil,
			expectedErrText: "Unknown command",
		},
		{
			name: "wrong symbols error",
			arg: "set k √",
			expectedCommand: "",
			expectedArgs: nil,
			expectedErrText: "Unknown symbols in arguments",
		},
	}

	for _, testCase := range testCases {
		actualCommand, actualArgs, actualErr := parser.ParseArgs(testCase.arg)
		if actualCommand != testCase.expectedCommand ||
		actualErr.Error() != testCase.expectedErrText {
			t.Errorf(
				`case %v:
				expectedCommand: %v 
				actualCommand: %v
				expectedArgs: %v
				actualArgs: %v
				expectedErrText: %v
				actualErrText: %v`,
				testCase.name,
				testCase.expectedCommand,
				actualCommand,
				testCase.expectedArgs,
				actualArgs,
				testCase.expectedErrText,
				actualErr.Error(),
			)
		}

		for i, expectedArg := range testCase.expectedArgs {
			if len(expectedArg) != len(actualArgs) || actualArgs[i] != expectedArg {
				t.Errorf(`case %v:
					expectedArgs: %v
					actualArgs: %v
				`,
				testCase.name,
				testCase.expectedArgs,
				actualArgs,
				)
			}
		}
	}
}
