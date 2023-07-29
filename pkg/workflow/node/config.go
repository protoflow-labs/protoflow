package node

//type ConfigNode struct {
//	BaseNode
//	Config *gen.Config
//}
//
//var _ graph.Node = &ConfigNode{}
//
//func NewConfigNode(node *gen.Node) *ConfigNode {
//	return &ConfigNode{
//		BaseNode: NewBaseNode(node),
//		Config:   node.GetConfiguration(),
//	}
//}
//
//func (c *ConfigNode) Represent() (string, error) {
//	return c.Config.Value, nil
//}
//
//func (c *ConfigNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
//	return graph.Output{}, errors.New("not implemented")
//}
