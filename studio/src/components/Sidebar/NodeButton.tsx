import {useEditorContext} from "@/providers/EditorProvider";
import {CardHeader, Text} from "@fluentui/react-components";
import {BaseBlockCard} from "@/components/blocks/BaseBlockCard";
import {Node as ProtoNode} from "@/rpc/graph_pb";
import React from "react";
import {BiClipboard} from "react-icons/bi";

export const NodeButton: React.FC<{ node: ProtoNode }> = ({ node }) => {
    const { setDraggedNode } = useEditorContext();
    return (
        <div
            className="m-2"
            style={{marginBottom: "10px"}}
            draggable
            onDragStart={(e) => {
                setDraggedNode(node);
            }}
        >
            <BaseBlockCard selected={false} style={{ cursor: "grab" }}>
                <CardHeader
                    image={<BiClipboard className="h-5 w-5 bg-gray-800" />}
                    header={<Text weight="semibold">{node.name}</Text>}
                />
                {node.name && <p style={{marginBottom: "0"}}>{node.name}</p>}
            </BaseBlockCard>
        </div>
    );
}
