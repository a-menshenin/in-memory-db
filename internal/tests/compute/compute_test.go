package compute

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"umemory/internal/compute"
	mock_compute "umemory/internal/compute/mock"
	mock_storage "umemory/internal/storage/mock"

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

	mockStorage := mock_storage.NewMockStorage(ctrl)
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
				mockParser.EXPECT().ParseArgs("set asd").Return("", []string{}, errors.New("ожидается 2 аргумента, получено 1"))
			},
			expected: "Ошибка при парсинге аргументов: ожидается 2 аргумента, получено 1",
		},
		{
			name: "Handle: get storage not found",
			requestStr: "get asd",
			exec: func() {
				mockParser.EXPECT().ParseArgs("get asd").Return("get", []string{"asd"}, nil)
				mockStorage.EXPECT().Get("asd").Return("", false)
			},
			expected: "Значение по ключу asd не найдено",
		},
		{
			name: "Handle: set value to storage",
			requestStr: "set key value",
			exec: func() {
				mockParser.EXPECT().ParseArgs("set key value").Return("set", []string{"key", "value"}, nil)
				mockStorage.EXPECT().Set("key", "value")
			},
			expected: "Значение value сохранено",
		},
		{
			name: "Handle: get value from storage",
			requestStr: "get key",
			exec: func() {
				mockParser.EXPECT().ParseArgs("get key").Return("get", []string{"key"}, nil)
				mockStorage.EXPECT().Get("key").Return("value", true)
			},
			expected: "Найдено значение: value",
		},
		{
			name: "Handle: delete value from storage",
			requestStr: "delete key",
			exec: func() {
				mockParser.EXPECT().ParseArgs("delete key").Return("delete", []string{"key"}, nil)
				mockStorage.EXPECT().Delete("key")
			},
			expected: "Значение key удалено",
		},
	}

	for _, testCase := range testCases {
		testCase.exec()
		getOutputText := captureFromStdout()
		handler.Handle(testCase.requestStr)
		actualOutput, err := getOutputText()
		if err != nil {
			t.Errorf("case %v: getOutputText error: %v \nactual output: %v", testCase.name, err, actualOutput)
		}
		if actualOutput != testCase.expected {
			t.Errorf("case %v: \nexpectedOutput: %v \nactualOutput: %v", testCase.name, testCase.expected, actualOutput)
		}
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
