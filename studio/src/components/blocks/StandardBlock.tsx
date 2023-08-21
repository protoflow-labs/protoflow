import {BaseBlockCard} from "./BaseBlockCard";
import {Caption1, CardHeader, Text} from "@fluentui/react-components";
import {Handle, NodeProps, Position} from "reactflow";
import {ReactNode} from "react";
import {useProjectContext} from "@/providers/ProjectProvider";
import {EnumeratedProvider, ProviderState } from "@/rpc/project_pb";
import { Edge as ProtoEdge, Node as ProtoNode } from "@/rpc/graph_pb";

export interface StandardBlockProps {
    name:string;
    description:string;
    image: ReactNode;
    selected: boolean;
}

type StandardBlockNodeProps = NodeProps<StandardBlockProps>;

function nodeIsOffline(node: ProtoNode, nodeLookup: Record<string, ProtoNode>, providerLookup: Record<string, EnumeratedProvider>, edgeLookup: Record<string, ProtoEdge>): boolean {
    const edges = Object.values(edgeLookup).filter(e => e.to === node.id && e.type.case === "provides");
    const provider = providerLookup[node.id];

    if (provider && provider.info) {
        return provider.info.state === ProviderState.ERROR;
    }
    return edges.some(e => nodeIsOffline(nodeLookup[e.from], nodeLookup, providerLookup, edgeLookup));
}

export function StandardBlock(props: StandardBlockNodeProps) {
    const { description, image, selected, id } = props;

    const {nodeLookup, providerLookup, edgeLookup} = useProjectContext();

    const node = nodeLookup[id];
    const isOffline = nodeIsOffline(node, nodeLookup, providerLookup, edgeLookup);

    return (
        <>
            <BaseBlockCard selected={selected} appearance={isOffline ? 'outline' : 'filled'}>
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
