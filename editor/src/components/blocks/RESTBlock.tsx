import { Caption1, CardHeader, Text } from "@fluentui/react-components";
import { MdHttp } from "react-icons/md";
import { Handle, NodeProps, Position } from "reactflow";
import { BlockCard } from "./BlockCard";

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
      <BlockCard selected={selected}>
        <CardHeader
          image={<MdHttp className="h-5 w-5 bg-gray-800" />}
          header={<Text weight="semibold">{data.name || "Untitled REST"}</Text>}
          description={<Caption1>REST</Caption1>}
        />
      </BlockCard>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}
