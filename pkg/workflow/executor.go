package workflow

import (
	"fmt"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
	"reflect"
	"runtime"
	"strings"
)

type Executor interface {
	Execute(activity interface{}, block interface{}, input Input) (*Result, error)
}

var _ Executor = &TemporalExecutor{}

type TemporalExecutor struct {
	ctx workflow.Context
}

func NewTemporalExecutor(ctx workflow.Context) *TemporalExecutor {
	return &TemporalExecutor{
		ctx: ctx,
	}
}

func (e *TemporalExecutor) Execute(activity interface{}, block interface{}, input Input) (*Result, error) {
	var result Result
	err := workflow.ExecuteActivity(e.ctx, activity, block, input).Get(e.ctx, &result)
	if err != nil {
		return nil, errors.Wrap(err, "error executing activity")
	}
	return &result, nil
}

type MemoryExecutor struct {
	ctx *MemoryContext
}

var _ Executor = &MemoryExecutor{}

func (e *MemoryExecutor) Execute(activity interface{}, block interface{}, input Input) (*Result, error) {
	activityArgs := []interface{}{
		e.ctx, block, input,
	}
	res, err := executeFunction(activity, activityArgs)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	result, ok := res.(Result)
	if !ok {
		return nil, fmt.Errorf("invalid result type: %T", res)
	}
	return &result, nil
}

func NewMemoryExecutor(ctx *MemoryContext) *MemoryExecutor {
	return &MemoryExecutor{
		ctx: ctx,
	}
}

func getFunctionName(i interface{}) (name string, isMethod bool) {
	if fullName, ok := i.(string); ok {
		return fullName, false
	}
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	// Full function name that has a struct pointer receiver has the following format
	// <prefix>.(*<type>).<function>
	isMethod = strings.ContainsAny(fullName, "*")
	elements := strings.Split(fullName, ".")
	shortName := elements[len(elements)-1]
	// This allows to call activities by method pointer
	// Compiler adds -fm suffix to a function name which has a receiver
	// Note that this works even if struct pointer used to get the function is nil
	// It is possible because nil receivers are allowed.
	// For example:
	// var a *Activities
	// ExecuteActivity(ctx, a.Foo)
	// will call this function which is going to return "Foo"
	return strings.TrimSuffix(shortName, "-fm"), isMethod
}

// Executes function and ensures that there is always 1 or 2 results and second
// result is error.
func executeFunction(fn interface{}, args []interface{}) (interface{}, error) {
	fnValue := reflect.ValueOf(fn)
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		// If the argument is nil, use zero value
		if arg == nil {
			reflectArgs[i] = reflect.New(fnValue.Type().In(i)).Elem()
		} else {
			reflectArgs[i] = reflect.ValueOf(arg)
		}
	}
	retValues := fnValue.Call(reflectArgs)

	// Expect either error or (result, error)
	if len(retValues) == 0 || len(retValues) > 2 {
		fnName, _ := getFunctionName(fn)
		return nil, fmt.Errorf(
			"the function: %v signature returns %d results, it is expecting to return either error or (result, error)",
			fnName, len(retValues))
	}
	// Convert error
	var err error
	if errResult := retValues[len(retValues)-1].Interface(); errResult != nil {
		var ok bool
		if err, ok = errResult.(error); !ok {
			return nil, fmt.Errorf(
				"failed to serialize error result as it is not of error interface: %v",
				errResult)
		}
	}
	// If there are two results, convert the first only if it's not a nil pointer
	var res interface{}
	if len(retValues) > 1 && (retValues[0].Kind() != reflect.Ptr || !retValues[0].IsNil()) {
		res = retValues[0].Interface()
	}
	return res, err
}