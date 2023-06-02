import {ReactNode} from "react";
import {Block} from "@/rpc/block_pb";
import {ReactFlowProtoflowData, ReactFlowProtoflowKey} from "@/providers/EditorProvider";
import { Button } from 'flowbite-react';
import {Caption1, CardHeader, Text} from "@fluentui/react-components";
import {HiPencilSquare} from "react-icons/hi2";
import {BaseBlockCard} from "@/components/blocks/BaseBlockCard";

export function NodeButton(props: {
    label: string,
    typeName: string,
    image: ReactNode,
    // The below optional params are used when dragging in blocks from existing services, such as protoflow itself, or an externally reflected GRPC service
    nodeName?: string,
    nodeConfig?: Block['type'],
    resourceIds?: string[],
    newBlock?: boolean,
}) {
    return (
        <div
            className="m-2"
            style={{marginBottom: "10px"}}
            draggable
            onDragStart={(e) => {
                const data: ReactFlowProtoflowData = {
                    type: 'protoflow.' + props.typeName,
                    config: props.nodeConfig,
                    name: props.nodeName,
                    resourceIds: props.resourceIds || [],
                }
                e.dataTransfer.setData(ReactFlowProtoflowKey, JSON.stringify(data));
                e.dataTransfer.effectAllowed = "move";
            }}
        >

            <BaseBlockCard selected={false} style={{ cursor: "grab" }}>
                <CardHeader
                    image={props.image}
                    header={<Text weight="semibold">{props.label}</Text>}
                />
            </BaseBlockCard>


            {/*<Button*/}
            {/*    size="medium"*/}
            {/*    className="w-full"*/}
            {/*    appearance={props.newBlock ? "primary" : "secondary"}*/}
            {/*    style={{ cursor: "grab" }}*/}
            {/*>*/}
            {/*    {props.children}*/}
            {/*</Button>*/}
        </div>
    );
}
