package storage

import (
	"testing"
	"umemory/internal/storage"
)

type storageTestCase struct {
	name string
	command string
	args []string
	expected string
	found bool
}

func TestInMemoryStorage(t *testing.T) {
	storage := storage.NewInMemoryStorage()

	var testCases = []storageTestCase{
		{
			name: "set value",
			command: "set",
			args: []string{"testkey", "testvalue"},
			expected: "",
			found: false,
		},
		{
			name: "get existing value",
			command: "get",
			args: []string{"testkey"},
			expected: "testvalue",
			found: true,
		},
		{
			name: "get not existing value",
			command: "get",
			args: []string{"notestkey"},
			expected: "",
			found: false,
		},
		{
			name: "delete value",
			command: "delete",
			args: []string{"testkey"},
			expected: "",
			found: false,
		},
		{
			name: "get deleted value",
			command: "get",
			args: []string{"testkey"},
			expected: "",
			found: false,
		},
	}

	for _, testCase := range testCases {
		switch testCase.command {
		case "set":
			storage.Set(testCase.args[0], testCase.args[1])
		case "get":
			actualRes, actualFound := storage.Get(testCase.args[0])
			if actualRes != testCase.expected || actualFound != testCase.found {
				t.Errorf("case %v: \nexpected value: %v \ngot value: %v \nexpected found: %v \ngot found: %v", testCase.name, testCase.expected, actualRes, testCase.found, actualFound)
			}
		case "delete":
			storage.Delete(testCase.args[0])
		}
	}
}
