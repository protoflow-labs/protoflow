import { Caption1, Card, CardHeader, Text } from "@fluentui/react-components";
import { MdHttp } from "react-icons/md";
import { Handle, NodeProps, Position } from "reactflow";

export type RESTBlockProps = NodeProps<RESTData>;

export type RESTData = {
  name: string;
  config: {
    rest: { method: string; path: string; headers: Record<string, string> };
  };
};

export function RESTBlock(props: RESTBlockProps) {
  const { data, selected } = props;

  return (
    <>
      <Card>
        <CardHeader
          image={<MdHttp className="h-5 w-5 bg-gray-800" />}
          header={<Text weight="semibold">{data.name || "Untitled REST"}</Text>}
          description={<Caption1>REST</Caption1>}
        />
      </Card>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}
