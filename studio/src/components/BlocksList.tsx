import {Accordion, AccordionHeader, AccordionItem, AccordionPanel, Button} from "@fluentui/react-components";
import { ReactNode } from "react";
import {useProjectResources} from "@/hooks/useProjectResources";
import {useProjectContext} from "@/providers/ProjectProvider";
import {ReactFlowProtoflowData, ReactFlowProtoflowKey} from "@/providers/EditorProvider";
import {Block} from "@/rpc/block_pb";

interface NodeBlock {
  type: string
  name: string
}
export default function BlocksList() {
  const { resources } = useProjectContext();
  const builtinBlocks: NodeBlock[] = [
      { type: "protoflow.input", name: "Input" },
      { type: "protoflow.collection", name: "Collection"},
      { type: "protoflow.function", name: "Function" },
      { type: "protoflow.query", name: "Query" },
      { type: "protoflow.queue", name: "Queue" },
      { type: "protoflow.bucket", name: "Bucket" },
      { type: "protoflow.email", name: "Email" },
      { type: "protoflow.rest", name: "REST" },
  ]
  return (
    <div className="absolute flex flex-col gap-1 m-3 z-10 top-8">
      <Accordion defaultOpenItems={"Built-in"}>
        <AccordionItem value="Built-in">
          <AccordionHeader>Built-in</AccordionHeader>
          <AccordionPanel>
            {builtinBlocks.map((block) => {
              return (
                <NodeButton key={block.type} nodeType={block.type}>
                  {block.name}
                </NodeButton>
              );
            })}
          </AccordionPanel>
        </AccordionItem>
        {resources.map((resource) => {
          return (
            <AccordionItem key={resource.id} value={resource.name}>
              <AccordionHeader>
                {resource.name}
              </AccordionHeader>
              <AccordionPanel>
                {resource.blocks.map((block) => {
                  // TODO breadchris support more block types
                  return (
                    <NodeButton
                      key={block.id}
                      nodeType="protoflow.grpc"
                      nodeName={block.name}
                      nodeConfig={block.type}
                      resourceIds={[resource.id]}
                    >
                      {block.name}
                    </NodeButton>
                  );
                })}
              </AccordionPanel>
            </AccordionItem>
          );
        })}
      </Accordion>
    </div>
  );
}

function NodeButton(props: {
  children: ReactNode,
  nodeType: string,
  nodeName?: string,
  nodeConfig?: Block['type'],
  resourceIds?: string[],
}) {
  return (
    <div
      draggable
      onDragStart={(e) => {
        const data: ReactFlowProtoflowData = {
          type: props.nodeType,
          config: props.nodeConfig,
          name: props.nodeName,
          resourceIds: props.resourceIds || [],
        }
        e.dataTransfer.setData(ReactFlowProtoflowKey, JSON.stringify(data));
        e.dataTransfer.effectAllowed = "move";
      }}
    >
      <Button size="small" className="w-full">
        {props.children}
      </Button>
    </div>
  );
}
