import {BaseBlockCard} from "./BaseBlockCard";
import {Caption1, CardHeader, Text} from "@fluentui/react-components";
import {Handle, NodeProps, Position} from "reactflow";
import {ReactNode} from "react";
import {useProjectContext} from "@/providers/ProjectProvider";

export interface StandardBlockProps {
    name:string;
    description:string;
    image: ReactNode;
    selected: boolean;
}

type StandardBlockNodeProps = NodeProps<StandardBlockProps>;


export function StandardBlock(props: StandardBlockNodeProps) {
    const { description, image, selected, id } = props;

    const {nodeLookup} = useProjectContext();

    const node = nodeLookup[id];
    return (
        <>
            <BaseBlockCard selected={selected}>
                <CardHeader
                    image={image}
                    header={
                        <Text weight="semibold">{node.name || "Untitled"}</Text>
                    }
                    description={<Caption1>{description}</Caption1>}
                />
            </BaseBlockCard>
            <Handle type="source" position={Position.Bottom} className="z-10" />
            <Handle type="target" position={Position.Top} className="z-10" />
        </>
    );
}
