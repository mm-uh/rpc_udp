package util

import "fmt"

func (CallerStruct) ExampleMethod(firstArg string) string {
	return "This is an example method with args: " + firstArg
}

func (CallerStruct) ExampleMethod2(firstArg string) string {
	return fmt.Sprintf("Example method2 with args %s", firstArg)
}
