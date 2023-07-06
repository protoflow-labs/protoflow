import { Card } from "@fluentui/react-components";
import {useOnSelectionChange} from "reactflow";
import { InputEditor } from "./blockEditors/InputEditor";
import NodeProvider from "@/providers/NodeProvider";
import {useProjectContext} from "@/providers/ProjectProvider";
import {Node as ProtoNode} from "@/rpc/graph_pb";
import {Bucket, Collection, Function, GRPC, Prompt, Query, REST} from "@/rpc/block_pb";
import {GenericNodeEditor} from "@/components/blockEditors/GenericNodeEditor";

export function EditorPanel() {
  const { activeNode, setActiveNodeId } = useProjectContext();

  useOnSelectionChange({
    onChange: ({ nodes }) => {
      if (nodes.length !== 1) {
        setActiveNodeId(null);
        return;
      }
      setActiveNodeId(nodes[0].id);
    },
  });

  if (!activeNode) {
    return null;
  }

  return (
    <NodeProvider nodeId={activeNode.id}>
      <div className="absolute top-0 right-0 m-4 z-10 overflow-auto">
        <Card>
          <NodeEditor node={activeNode} />
        </Card>
      </div>
    </NodeProvider>
  );
}

type NodeEditorProps = {
  node: ProtoNode | null;
};

function NodeEditor(props: NodeEditorProps) {
  if (!props.node) {
    return null;
  }

  switch (props.node.config.case) {
    case "input":
      return <InputEditor node={props.node} />;
    case "collection":
      return <GenericNodeEditor node={props.node} nodeConfig={'collection'} nodeConfigType={Collection} />
    case "query":
      return <GenericNodeEditor node={props.node} nodeConfig={'query'} nodeConfigType={Query} />
    case "function":
      return <GenericNodeEditor node={props.node} nodeConfig={"function"} nodeConfigType={Function} />;
    case "bucket":
      return <GenericNodeEditor node={props.node} nodeConfig={"bucket"} nodeConfigType={Bucket} />;
    case "rest":
      return <GenericNodeEditor node={props.node} nodeConfig={"rest"} nodeConfigType={REST} />;
    case "grpc":
      return <GenericNodeEditor node={props.node} nodeConfig={"grpc"} nodeConfigType={GRPC} />;
    case "prompt":
      return <GenericNodeEditor node={props.node} nodeConfig={"prompt"} nodeConfigType={Prompt} />;
    default:
      return null;
  }
}
