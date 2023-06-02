import { useProjectContext } from "@/providers/ProjectProvider";
import {Block, GRPC, Function} from "@/rpc/block_pb";
import {NodeButton} from "@/components/newBlockSidebar/NodeButton";
import { Button, Accordion } from 'flowbite-react';
import {blockTypes} from "@/components/blocks/blockTypes";
import {HiCodeBracket} from "react-icons/hi2";
import React from "react";


export default function BlocksList() {
  const { resources, deleteResource } = useProjectContext();

  const buildInBlockNames = ['input','collection','query','bucket','rest'];

  const defaultFunctionConfig: Block['type'] = {
    case: 'function',
    value: new Function({
      runtime: "node",
      grpc: new GRPC({
        package: "protoflow",
        service: "nodeService",
        method: "Method",
      }),
    }),
  };
  return (
    <div className="absolute flex flex-col gap-1 m-3 z-10 top-8 pt-5" style={{marginTop:"40px", marginLeft:"10px"}}>
      <Accordion >
        <Accordion.Panel value="Built-in">
          <Accordion.Title className="text-white-500">Standard Blocks</Accordion.Title>
          <Accordion.Content>
            {buildInBlockNames.map((typeName) => {
              // use the block types from the "blocks" folder since those dictate what icons and labels to use for a given block
              const blockConstants = blockTypes.find((block) => block.typeName === typeName);
              if (!blockConstants) {
                throw new Error("failed to lookup block information for the name " + typeName)
              }
              return (
                <NodeButton key={blockConstants.typeName} nodeType={"protoflow.bucket"} label={blockConstants.label} typeName={blockConstants.typeName} image={blockConstants.image}  >
                </NodeButton>
              );
            })}
          </Accordion.Content>
        </Accordion.Panel>
        {resources.map((resource) => {
          // TODO: can we get this to a point where its all in one node loop instead of all this forked logic?
          return (
            <Accordion.Panel key={resource.id} value={resource.name}>
              <Accordion.Title>
                {resource.name.charAt(0).toUpperCase() + resource.name.slice(1)}
              </Accordion.Title>
              <Accordion.Content>
                {resource.blocks.length === 0 && resource.type.case !== "languageService" && (
                  <div className="text-gray-400">No blocks</div>
                )}
                {resource.type.case === 'languageService' && (
                    <NodeButton
                        image={<HiCodeBracket className="h-5 w-5 bg-gray-800" />}
                        typeName={'function'}
                        nodeName={'new'}
                        label={"New"}

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
                      label={block.name}
                      typeName={block.type.case || "grpc"}
                      nodeConfig={block.type}
                      resourceIds={[resource.id]}
                    >
                      {block.name}
                    </NodeButton>
                  );
                })}
                <Button size="small" className="w-full" appearance={'outline'} onClick={() => deleteResource(resource.id)}>Delete</Button>
              </Accordion.Content>
            </Accordion.Panel>
          );
        })}
      </Accordion>
    </div>
  );
}


