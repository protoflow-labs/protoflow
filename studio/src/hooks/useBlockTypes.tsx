import { BucketBlock } from "@/components/blocks/BucketBlock";
import { CollectionBlock } from "@/components/blocks/CollectionBlock";
import { FunctionBlock } from "@/components/blocks/FunctionBlock";
import { InputBlock } from "@/components/blocks/InputBlock";
import { QueryBlock } from "@/components/blocks/QueryBlock";
import { QueueBlock } from "@/components/blocks/QueueBlock";
import { RESTBlock } from "@/components/blocks/RESTBlock";
import { useMemo } from "react";
import {GRPCBlock} from "@/components/blocks/GRPCBlock";
import {blockTypes} from "@/components/blocks/blockTypes";
import {StandardBlock} from "@/components/blocks/StandardBlock";


/*
export interface StandardBlockProps {
    name:string;
    description:string;
    image: ReactNode;
    selected: boolean;
}

 */

const nodeElementMap: Record<string, ReactNode> = {}
 blockTypes.forEach((blockType) => {

    nodeElementMap['protoflow.' + blockType.typeName] = (props:StandardBlockProps) => {
        return (
            <StandardBlock image={blockType.image} description={blockType.label} {...props} />
        );
    };
})

export function useBlockTypes() {

  return { nodeTypes: nodeElementMap };
}
