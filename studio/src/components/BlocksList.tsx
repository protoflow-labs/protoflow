import { Accordion, AccordionHeader, AccordionItem, AccordionPanel, Button } from "@fluentui/react-components";
import { ReactNode } from "react";
import { useProjectContext } from "@/providers/ProjectProvider";
import { ReactFlowProtoflowData, ReactFlowProtoflowKey } from "@/providers/EditorProvider";
import {Block, GRPC, Function} from "@/rpc/block_pb";

interface NodeBlock {
  type: string
  name: string
}
export default function BlocksList() {
  const { resources, deleteResource } = useProjectContext();
  const builtinBlocks: NodeBlock[] = [
    { type: "protoflow.input", name: "Input" },
    { type: "protoflow.collection", name: "Collection" },
    { type: "protoflow.query", name: "Query" },
    { type: "protoflow.bucket", name: "Bucket" },
    { type: "protoflow.rest", name: "REST" },
  ]
  const defaultFunctionConfig: Block['type'] = {
    case: 'function',
    value: new Function({
      runtime: "nodejs",
      grpc: new GRPC({
        package: "protoflow",
        service: "nodejsService",
        method: "Method",
      }),
    }),
  };
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
                {resource.blocks.length === 0 && (
                  <div className="text-gray-400">No blocks</div>
                )}
                {resource.type.case === 'languageService' && (
                    <NodeButton
                        nodeType={'protoflow.function'}
                        nodeName={'new'}
                        nodeConfig={defaultFunctionConfig}
                        resourceIds={[resource.id]}
                        newBlock={true}
                    >New</NodeButton>
                )}
                {resource.blocks.map((block) => {
                  // TODO breadchris support more block types
                  return (
                    <NodeButton
                      key={block.id}
                      nodeType={block.type.case === 'grpc' ? 'protoflow.grpc' : 'protoflow.function'}
                      nodeName={block.name}
                      nodeConfig={block.type}
                      resourceIds={[resource.id]}
                    >
                      {block.name}
                    </NodeButton>
                  );
                })}
                <Button size="small" className="w-full" appearance={'outline'} onClick={() => deleteResource(resource.id)}>Delete</Button>
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
  newBlock?: boolean,
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
      <Button size="small" className="w-full" appearance={props.newBlock ? "primary" : "secondary"}>
        {props.children}
      </Button>
    </div>
  );
}
