import {Block} from "@/rpc/block_pb";
import {ReactFlowProtoflowData, ReactFlowProtoflowKey, useEditorContext} from "@/providers/EditorProvider";
import { Button } from 'flowbite-react';
import {Caption1, CardHeader, Text} from "@fluentui/react-components";
import {HiPencilSquare} from "react-icons/hi2";
import {BaseBlockCard} from "@/components/blocks/BaseBlockCard";
import {blockTypes} from "@/components/blocks/blockTypes";
import {Node as ProtoNode} from "@/rpc/graph_pb";

export function NodeButton(props: {

    // The below optional params are used when dragging in blocks from existing services, such as protoflow itself, or an externally reflected GRPC service
    node: ProtoNode
}) {
    console.log('nodeButton rendered with props ', props);
    const { setDraggedNode } = useEditorContext();
    const blockStaticInfo = blockTypes.find((block) => block.typeName === props.node.config.case);
    console.log("static info about block is ", blockStaticInfo)
    return (
        <div
            className="m-2"
            style={{marginBottom: "10px"}}
            draggable
            onDragStart={(e) => {
                console.log("drag start", props.node);
                setDraggedNode(props.node);
            }}
        >

            <BaseBlockCard selected={false} style={{ cursor: "grab" }}>
                <CardHeader
                    image={blockStaticInfo.image}
                    header={<Text weight="semibold">{blockStaticInfo.label}</Text>}
                />
                {props.node.name && <p>{props.node.name}</p>}
            </BaseBlockCard>
        </div>
    );
}
