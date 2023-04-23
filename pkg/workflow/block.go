package workflow

import (
	"github.com/pkg/errors"
	protoflow "github.com/protoflow-labs/protoflow/gen"
)

type Block interface {
	Execute(executor Executor, input Input) (*Result, error)
}

type Code struct {
	Code string
}

type Data struct {
	Params interface{}
}

var activity = &Activity{}

func (s *Code) Init() error {
	return nil
}

func (s *Code) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteCode, s, input)
}

func (s *Data) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteInput, s, input)
}

func NewBlock(node *protoflow.Node) (Block, error) {
	switch node.Type.(type) {
	case *protoflow.Node_Function:
		f := node.GetFunction()
		switch f.Type.(type) {
		case *protoflow.Function_Code:
			c := f.GetCode()
			return &Code{
				Code: c.Code,
			}, nil
		default:
			return nil, errors.New("no code found")
		}
	case *protoflow.Node_Data:
		d := node.GetData()
		switch d.Type.(type) {
		case *protoflow.Data_Input:
			i := d.GetInput()
			return &Data{
				Params: i.Params,
			}, nil
		}
	default:
		return nil, errors.New("no function or data found")
	}
	return nil, nil
}
