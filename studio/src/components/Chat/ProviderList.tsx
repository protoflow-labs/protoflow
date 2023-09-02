import {
    Accordion,
    AccordionHeader,
    AccordionItem,
    AccordionPanel, CardHeader, Text,
} from "@fluentui/react-components";
import { useProjectContext } from "@/providers/ProjectProvider";
import { ProviderState } from "@/rpc/project_pb";
import { PlugDisconnected20Regular } from "@fluentui/react-icons";
import React from "react";
import {useEditorContext} from "@/providers/EditorProvider";
import {BaseNodeCard} from "@/components/BaseNodeCard";
import {BiClipboard} from "react-icons/bi";
import {Node as ProtoNode, NodeDetails} from "@/rpc/graph_pb";

const NodeButton: React.FC<{ provider: NodeDetails, node: ProtoNode }> = ({ node, provider }) => {
    const { setDraggedNode } = useEditorContext();
    return (
        <div
            className="m-2"
            style={{marginBottom: "10px"}}
            draggable
            onDragStart={(e) => {
                setDraggedNode({node, provider});
            }}
        >
            <BaseNodeCard selected={false} style={{ cursor: "grab" }}>
                <CardHeader
                    image={<BiClipboard className="h-5 w-5 bg-gray-800" />}
                    header={<Text weight="semibold">{node.name}</Text>}
                />
                {node.name && <p style={{marginBottom: "0"}}>{node.name}</p>}
            </BaseNodeCard>
        </div>
    );
}

export default function ProviderList() {
  const { providers } = useProjectContext();

  return (
      <div>
        <Accordion defaultOpenItems={"Built-in"} collapsible={true}>
          {providers.map((r) => {
            const p = r.provider;
            if (!p) {
                return null;
            }
            const resError = r.info && r.info.state === ProviderState.ERROR;
            return (
                <AccordionItem key={p.id} value={p.name}>
                  <AccordionHeader icon={resError ? <PlugDisconnected20Regular /> : null}>
                    {p.name}
                  </AccordionHeader>
                  <AccordionPanel className={"overflow-y-auto"} style={{maxHeight: "40em"}}>
                    {resError ? <p className={"text-red-500"}>Error: {r.info?.error}</p> : null}
                    {r.nodes.map((node) => {
                        return (
                            <NodeButton key={node.id} provider={p} node={node} />
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

