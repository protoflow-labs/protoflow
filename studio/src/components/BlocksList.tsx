import { Accordion, AccordionHeader, AccordionItem, AccordionPanel, Button } from "@fluentui/react-components";
import { ReactNode } from "react";
import { useProjectContext } from "@/providers/ProjectProvider";
import {ReactFlowProtoflowData, ReactFlowProtoflowKey, useEditorContext} from "@/providers/EditorProvider";
import {GRPC, Function, Input, Collection, Bucket} from "@/rpc/block_pb";
import {Node} from "@/rpc/graph_pb";
import {Resource} from "@/rpc/resource_pb";
import {generateUUID} from "@/util/uuid";

interface NodeBlock {
  type: string
  name: string
}

function resourceToNode(res: Resource, name: string) {
  const baseNode = new Node({
    id: generateUUID(),
    name: name,
    resourceId: res.id,
  });

  let config: Node['config'] | null = null;
  switch (res.type.case) {
    case 'languageService':
      config = {
        case: 'function',
        value: new Function({
          runtime: res.type.value.runtime,
        })
      }
      break;
    case 'docstore':
      config = {
          case: 'collection',
          value: new Collection({
          })
      }
      break;
    case 'blobstore':
      config = {
        case: 'bucket',
        value: new Bucket({})
      }
      break;
  }
  if (config !== null) {
    baseNode.config = config;
  }
  return baseNode;
}

export default function BlocksList() {
  const { resources, deleteResource } = useProjectContext();
  const builtinBlocks: Node[] = [
    new Node({
      id: generateUUID(),
      name: "Input",
      config: {
        case: "input",
        value: new Input({})
      }
    })
  ];

  return (
    <div className="absolute flex flex-col gap-1 m-3 z-10 top-8">
      <Accordion defaultOpenItems={"Built-in"}>
        <AccordionItem value="Built-in">
          <AccordionHeader>Built-in</AccordionHeader>
          <AccordionPanel>
            {builtinBlocks.map((node, i) => {
              return (
                <NodeButton key={i} node={node}>
                  {node.name}
                </NodeButton>
              );
            })}
          </AccordionPanel>
        </AccordionItem>
        {resources.map((r) => {
          if (!r.resource || !r.resource.type || !r.resource.type.case) {
            return null;
          }
          const res = r.resource;
          const n = resourceToNode(res, "new");
          return (
            <AccordionItem key={res.id} value={res.name}>
              <AccordionHeader>
                {res.name}
              </AccordionHeader>
              <AccordionPanel>
                {r.nodes.length === 0 && (
                  <div className="text-gray-400">No nodes</div>
                )}
                {n && (
                    <NodeButton node={n} newBlock={true}>New</NodeButton>
                )}
                {r.nodes.map((node) => {
                  // TODO breadchris support more block types
                  return (
                    <NodeButton key={node.id} node={node}>
                      {node.name}
                    </NodeButton>
                  );
                })}
                <Button size="small" className="w-full" appearance={'outline'} onClick={() => deleteResource(res.id)}>Delete</Button>
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

  node: Node,
  newBlock?: boolean,
}) {
  const { props: {onDrop }, setDraggedNode } = useEditorContext();
  return (
    <div
      draggable
      onDragStart={(e) => {
        console.log("drag start", props.node);
        setDraggedNode(props.node);
      }}
    >
      <Button size="small" className="w-full" appearance={props.newBlock ? "primary" : "secondary"}>
        {props.children}
      </Button>
    </div>
  );
}
