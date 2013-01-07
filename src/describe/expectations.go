package describe

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

type nilValueType interface{}
var nilValue *nilValueType = nil

func appendValueFor(array []reflect.Value, obj interface{}) []reflect.Value {
	var value reflect.Value
	if reflect.TypeOf(obj) == nil {
		value = reflect.ValueOf(nilValue)
	} else {
		value = reflect.ValueOf(obj)
	}
	return append(array, value)
}

// A matcher generator that negates the given matcher.
// Unlike other matchers, this function directly accepts the matcher in question:
//
// Example:
//
//    Expect(t, "red", Not(ToEqual), "blue")
//
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

// A matcher that expects the value to be true
func ToBeTrue(actual interface{}) (string, bool) {
	return ToEqual(actual, true)
}

// A matcher that expects the value to be false
func ToBeFalse(actual interface{}) (string, bool) {
    return ToEqual(actual, false)
}

// A matcher that expects the value to be nil
func ToBeNil(actual interface{}) (string, bool) {
	value := reflect.ValueOf(actual)
	if value.Kind() != reflect.Ptr || !value.IsNil() {
		return "to be nil", false
	}
	return "", true
}

// Expects the given value to have a length of the provided value
func ToBeLengthOf(actual interface{}, size int) (string, bool) {
	value := reflect.ValueOf(actual)
	if value.Len() != size {
		return fmt.Sprintf("to be length of %d (got %d)", size, value.Len()), false
	}
	return "", true
}

// Expects the given value to be have a length of zero
func ToBeEmpty(actual interface{}) (string, bool) {
	value := reflect.ValueOf(actual)
	if value.Len() != 0 {
		return fmt.Sprintf("to be empty (got %d)", value.Len()), false
	}
	return "", true
}

// Performs a simple equality comparison. Does not perform a deep equality.
func ToBe(actual, expected interface{}) (string, bool) {
	if actual != expected {
		return fmt.Sprintf("to equal %#v", expected), false
	}
	return "", true
}

// Performs a deep equal - comparing struct, arrays, and slices items too.
func ToEqual(actual, expected interface{}) (string, bool) {
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

// The interface that Expect requires for its first argument. It is used to report
// the results of the expectation.
type Reporter interface {
	FailNow()
	Logf(format string, args ...interface{})
}

// Performs an expectation. It takes a Reporter (which testing.T satisfies), followed
// by the value under test, then a matcher. Any additional arguments, after that
// are passed directly to the matcher. Certain matches may require more arguments.
//
// Example:
//
//    func TestBoolean(t *testing.T) {
//      Expect(t, true, ToBeTrue)
//      Expect(t, 1, ToEqual, 1)
//    }
//
func Expect(r Reporter, obj interface{}, test interface{}, args ...interface{}) {
	var argValues []reflect.Value

	argValues = appendValueFor(argValues, obj)
	for _, v := range args {
		argValues = appendValueFor(argValues, v)
	}

	testfn := reflect.ValueOf(test)
	if testfn.Kind() != reflect.Func {
		stacktrace := tabulate("Stacktrace: ", getStackTrace(3), "\n")
        r.Logf("Expect() requires 3rd argument to be matcher func\n       got: %#v\n\n%s", test, stacktrace)
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

// Fails the test immediately
func Fail(r Reporter, message string) {
    stacktrace := tabulate(" stacktrace: ", getStackTrace(3), "\n")
    r.Logf("Fail %s:\n%s", message, stacktrace)
    r.FailNow()
}
