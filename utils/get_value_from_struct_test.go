package utils_test

import (
	"github.com/wesdean/story-book-api/utils"
	"testing"
)

func TestGetValueFromStruct(t *testing.T) {
	type MyNestedStruct struct {
		MyNestedValue string
	}

	type MyStruct struct {
		Value        string
		NestedStruct MyNestedStruct
	}

	myStruct := MyStruct{Value: "t Value", NestedStruct: MyNestedStruct{MyNestedValue: "t NestedValue"}}
	myValue, err := utils.GetValueFromStruct("MyStruct.NestedStruct.MyNestedValue", &myStruct)
	if err != nil {
		t.Error(err)
		return
	}

	myString, ok := myValue.(string)
	if !ok {
		t.Errorf("expected string, got %T", myValue)
		return
	}

	expected := "t NestedValue"
	if myString != expected {
		t.Errorf("expected %v, got %v", expected, myValue)
		return
	}
}
