package describe

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

type validNilType interface{}

func validNil() *validNilType { return nil }

func appendValueFor(array []reflect.Value, obj interface{}) []reflect.Value {
	var value reflect.Value
	if reflect.TypeOf(obj) == nil {
		value = reflect.ValueOf(validNil())
	} else {
		value = reflect.ValueOf(obj)
	}
	return append(array, value)
}

func Not(test interface{}) func(actual interface{}, args ...interface{}) (string, bool) {
	return func(actual interface{}, args ...interface{}) (string, bool) {
		var argValues []reflect.Value

		argValues = appendValueFor(argValues, actual)
		for _, v := range args {
			argValues = appendValueFor(argValues, v)
		}

		fmt.Printf("Func: %#v; Args: %#v", test, argValues)
		returnValues := reflect.ValueOf(test).Call(argValues)
		str, ok := returnValues[0].String(), returnValues[1].Bool()

		return fmt.Sprintf("not %s (%v)", str), !ok

	}
}

func ToBeTrue(actual interface{}) (string, bool) {
	return ToEqual(actual, true)
}

func ToBeNil(actual interface{}) (string, bool) {
	value := reflect.ValueOf(actual)
	if !value.IsNil() {
		return "to be nil", false
	}
	return "", true
}

func ToBeLengthOf(actual interface{}, size int) (string, bool) {
	value := reflect.ValueOf(actual)
	if value.Len() != size {
		return fmt.Sprintf("to be length of %d (got %d)", size, value.Len()), false
	}
	return "", true
}

func ToBeEmpty(actual interface{}) (string, bool) {
	value := reflect.ValueOf(actual)
	if value.Len() != 0 {
		return fmt.Sprintf("to be empty (got %d)", value.Len()), false
	}
	return "", true
}

func ToFail(actual interface{}, message string) (string, bool) {
	return message, false
}

func ToEqual(actual, expected interface{}) (string, bool) {
	if actual != expected {
		return fmt.Sprintf("to equal %#v", expected), false
	}
	return "", true
}

func ToDeeplyEqual(actual, expected interface{}) (string, bool) {
	if !reflect.DeepEqual(actual, expected) {
		return fmt.Sprintf("to deeply equal %#v", expected), false
	}
	return "", true
}

func getStackTrace(offset int) string {
	strstack := make([]string, 0)
	stack := make([]uintptr, 100)
	count := runtime.Callers(offset, stack)
	stack = stack[0:count]
	for i, pc := range stack {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		if i == 0 {
			line--
		}
		strstack = append(strstack, fmt.Sprintf("%s:%d\n    inside %s", file, line, fn.Name()))
	}
	return strings.Join(strstack, "\n")
}

func tabulate(prefix, content, sep string) string {
	lines := strings.Split(content, sep)
	buffer := strings.Repeat(" ", len(prefix))
	for i, line := range lines {
		if i != 0 {
			lines[i] = buffer + line
		}
	}
	return prefix + strings.Join(lines, sep)
}

type Reporter interface {
	FailNow()
	Logf(format string, args ...interface{})
}

func Expect(r Reporter, obj interface{}, test interface{}, args ...interface{}) {
	var argValues []reflect.Value

	argValues = appendValueFor(argValues, obj)
	for _, v := range args {
		argValues = appendValueFor(argValues, v)
	}

	testfn := reflect.ValueOf(test)
	if testfn.Kind() != reflect.Func {
		stacktrace := tabulate("Stacktrace: ", getStackTrace(3), "\n")
		fmt.Fprintf(os.Stderr, "Expect() requires 3rd argument to be matcher func\n       got: %#v\n\n%s", test, stacktrace)
		os.Exit(1)
	}

	returnValues := testfn.Call(argValues)
	str, ok := returnValues[0].String(), returnValues[1].Bool()
	if !ok {
		stacktrace := tabulate(" stacktrace: ", getStackTrace(3), "\n")
		r.Logf("expected %#v %s\n%s", obj, str, stacktrace)
		r.FailNow()
	}
}
