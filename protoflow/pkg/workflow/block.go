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
	Code string
}

var activity = &Activity{}

func (s *Code) Init() error {
	return nil
}

func (s *Code) Execute(ctx workflow.Context) (interface{}, error) {
	res := Result{}
	err := workflow.ExecuteActivity(ctx, activity.ExecuteCode, s).Get(ctx, &res)
	return res, err
}

type Input struct {
	Params map[string]string
}

func (s *Input) Execute(ctx workflow.Context) (interface{}, error) {
	res := Result{}
	err := workflow.ExecuteActivity(ctx, activity.ExecuteInput, s).Get(ctx, &res)
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
			return &Input{
				Params: i.Params,
			}, nil
		}
	default:
		return nil, errors.New("no function or data found")
	}
	return nil, nil
}
