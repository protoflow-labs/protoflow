package workflow

import (
	protoflow "github.com/breadchris/protoflow/gen/workflow"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
)

type Block interface {
	Execute(ctx workflow.Context) (interface{}, error)
}

type Code struct {
	code string
}

func (s *Code) Execute(ctx workflow.Context) (interface{}, error) {
	res := Result{}
	err := workflow.ExecuteActivity(ctx, Activity.ExecuteCode, s).Get(ctx, &res)
	return res, err
}

type Input struct {
	params map[string]string
}

func (s *Input) Execute(ctx workflow.Context) (interface{}, error) {
	res := Result{}
	err := workflow.ExecuteActivity(ctx, Activity.ExecuteInput, s).Get(ctx, &res)
	return res, err
}

func NewBlock(node *protoflow.Node) (Block, error) {
	switch node.Type.(type) {
	case *protoflow.Node_Function:
		f := node.GetFunction()
		switch f.Type.(type) {
		case *protoflow.Function_Code:
			c := f.GetCode()
			return &Code{
				code: c.Code,
			}, nil
		default:
			return nil, errors.New("no code found")
		}
	case *protoflow.Node_Data:
		d := node.GetData()
		switch d.Type.(type) {
		case *protoflow.Data_Input:
			i := d.GetInput()
			return &Input{
				params: i.Params,
			}, nil
		}
	default:
		return nil, errors.New("no function or data found")
	}
	return nil, nil
}
