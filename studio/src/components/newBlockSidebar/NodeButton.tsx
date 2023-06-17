import {useEditorContext} from "@/providers/EditorProvider";
import {CardHeader, Text} from "@fluentui/react-components";
import {BaseBlockCard} from "@/components/blocks/BaseBlockCard";
import {blockTypes} from "@/components/blocks/blockTypes";
import {Node as ProtoNode} from "@/rpc/graph_pb";
import React from "react";

export const NodeButton: React.FC<{ node: ProtoNode }> = ({ node }) => {
    const { setDraggedNode } = useEditorContext();
    const blockStaticInfo = blockTypes.find((block) => block.typeName === node.config.case);
    if (!blockStaticInfo) {
        return null;
    }
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
                    image={blockStaticInfo.image}
                    header={<Text weight="semibold">{blockStaticInfo.label}</Text>}
                />
                {node.name && <p style={{marginBottom: "0"}}>{node.name}</p>}
            </BaseBlockCard>
        </div>
    );
}
