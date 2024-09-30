package compute

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"umemory/internal/compute"
	mock_compute "umemory/internal/compute/mock"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

type computeTestCase struct {
	name string
	requestStr string
	exec func()
	expected string
}

func TestComputeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_compute.NewMockStorage(ctrl)
	mockParser := mock_compute.NewMockParser(ctrl)

	handler := compute.NewComputeHandler(
		mockStorage,
		mockParser,
		zap.NewNop(),
	)

	var testCases = []computeTestCase{
		{
			name: "Handle: parser error",
			requestStr: "set asd",
			exec: func() {
				mockParser.EXPECT().ParseArgs("set asd").Return("", []string{}, errors.New("expected 2 arguments, got 1"))
			},
			expected: "Arguments parse error: expected 2 arguments, got 1",
		},
		{
			name: "Handle: get storage not found",
			requestStr: "get asd",
			exec: func() {
				mockParser.EXPECT().ParseArgs("get asd").Return("get", []string{"asd"}, nil)
				mockStorage.EXPECT().Get("asd").Return("", false)
			},
			expected: "Value by key=asd not found\n",
		},
		{
			name: "Handle: set value to storage",
			requestStr: "set key value",
			exec: func() {
				mockParser.EXPECT().ParseArgs("set key value").Return("set", []string{"key", "value"}, nil)
				mockStorage.EXPECT().Set("key", "value")
			},
			expected: "Value value saved\n",
		},
		{
			name: "Handle: get value from storage",
			requestStr: "get key",
			exec: func() {
				mockParser.EXPECT().ParseArgs("get key").Return("get", []string{"key"}, nil)
				mockStorage.EXPECT().Get("key").Return("value", true)
			},
			expected: "Value found: value\n",
		},
		{
			name: "Handle: delete value from storage",
			requestStr: "delete key",
			exec: func() {
				mockParser.EXPECT().ParseArgs("delete key").Return("delete", []string{"key"}, nil)
				mockStorage.EXPECT().Delete("key")
			},
			expected: "Value key deleted\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.exec()
			getOutputText := captureFromStdout()
			handler.Handle(tt.requestStr)
			actualOutput, err := getOutputText()
			if err != nil {
				t.Errorf("getOutputText error: %v \nactual output: %v", err, actualOutput)
			}
			if actualOutput != tt.expected {
				t.Errorf("expectedOutput: %v \nactualOutput: %v", tt.expected, actualOutput)
			}
		})
	}
}

func captureFromStdout() func() (string, error) {
    r, w, err := os.Pipe()
    if err != nil {
        panic(err)
    }

    done := make(chan error, 1)

    save := os.Stdout
    os.Stdout = w

    var buf strings.Builder

    go func() {
        _, err := io.Copy(&buf, r)
        r.Close()
        done <- err
    }()

    return func() (string, error) {
        os.Stdout = save
        w.Close()
        err := <-done
        return buf.String(), err
    }
}
