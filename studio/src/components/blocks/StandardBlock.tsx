import {BaseBlockCard} from "./BaseBlockCard";
import {Caption1, CardHeader, Text} from "@fluentui/react-components";
import {Handle, NodeProps, Position} from "reactflow";
import {ReactNode} from "react";

export interface StandardBlockProps {
    name:string;
    description:string;
    image: ReactNode;
    selected: boolean;
}

type StandardBlockNodeProps = NodeProps<StandardBlockProps>;


export function StandardBlock(props: StandardBlockNodeProps) {
    const { data, description, image, selected } = props;

    return (
        <>
            <BaseBlockCard selected={selected}>
                <CardHeader
                    image={image}
                    header={
                        <Text weight="semibold">{data.name || "Untitled"}</Text>
                    }
                    description={<Caption1>{description}</Caption1>}
                />
            </BaseBlockCard>
            <Handle type="source" position={Position.Bottom} className="z-10" />
            <Handle type="target" position={Position.Top} className="z-10" />
        </>
    );
}
