import {
  Accordion,
  AccordionHeader,
  AccordionItem,
  AccordionPanel,
  Button,
  Popover, Select,
  Tooltip
} from "@fluentui/react-components";
import { useProjectContext } from "@/providers/ProjectProvider";
import {GRPC, Function, Input, Collection, Bucket, Prompt, Config, Query, Template, Route} from "@/rpc/block_pb";
import {Node} from "@/rpc/graph_pb";
import {Resource} from "@/rpc/resource_pb";
import {NodeButton} from "@/components/Sidebar/NodeButton";
import { ProviderState } from "@/rpc/project_pb";
import { PlugDisconnected20Regular } from "@fluentui/react-icons";
import {useState} from "react";
import { Data } from "@/rpc/data/data_pb";
import {StandardBlock} from "@/components/blocks/StandardBlock";

interface NodeBlock {
  type: string
  name: string
}

export default function BlocksList() {
  const { providers } = useProjectContext();

  return (
      <div className="absolute flex flex-col gap-1 m-3 z-10 top-8" style={{marginTop: "40px"}}>
        <Accordion defaultOpenItems={"Built-in"} collapsible={true}>
          {providers.map((r) => {
            if (!r.provider) {
                return null;
            }
            const resError = r.info && r.info.state === ProviderState.ERROR;
            const res = r.provider;
            return (
                <AccordionItem key={res.id} value={res.name} disabled={resError}>
                  <AccordionHeader icon={resError ? <PlugDisconnected20Regular /> : null}>
                    {res.name}
                  </AccordionHeader>
                  <AccordionPanel className={"overflow-y-auto"} style={{maxHeight: "40em"}}>
                    {r.nodes.map((node) => {
                        return (
                            <NodeButton key={node.id} node={node} />
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

