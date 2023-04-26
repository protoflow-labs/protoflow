import { Caption1, Card, CardHeader, Text } from "@fluentui/react-components";
import { HiEnvelope } from "react-icons/hi2";
import { Handle, NodeProps, Position } from "reactflow";

export type EmailBlockProps = NodeProps<EmailData>;

export type EmailData = {
  name: string;
};

export function EmailBlock(props: EmailBlockProps) {
  const { data, selected } = props;
  return (
    <>
      <Card>
        <CardHeader
          image={<HiEnvelope className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{data.name || "On Termination Email"}</Text>
          }
          description={<Caption1>Email Trigger</Caption1>}
        />
      </Card>
      <Handle type="source" position={Position.Bottom} className="z-10" />
    </>
  );
}
