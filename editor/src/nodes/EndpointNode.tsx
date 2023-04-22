import { Caption1, Card, CardHeader, Text } from "@fluentui/react-components";
import { TbPlug } from "react-icons/tb";
import { Handle, NodeProps, Position } from "reactflow";

export type EndpointyNodeProps = NodeProps<EndpointyData>;

export type EndpointyData = {
  name: string;
};

export function EndpointyNode(props: EndpointyNodeProps) {
  const { data, selected } = props;
  return (
    <>
      <Card>
        <CardHeader
          image={<TbPlug className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{data.name || "Untitled Endpoint"}</Text>
          }
          description={<Caption1>Endpoint</Caption1>}
        />
      </Card>
      <Handle type="source" position={Position.Bottom} className="z-10" />
    </>
  );
}
