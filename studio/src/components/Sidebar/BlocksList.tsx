import {
  Accordion,
  AccordionHeader,
  AccordionItem,
  AccordionPanel,
} from "@fluentui/react-components";
import { useProjectContext } from "@/providers/ProjectProvider";
import {NodeButton} from "@/components/Sidebar/NodeButton";
import { ProviderState } from "@/rpc/project_pb";
import { PlugDisconnected20Regular } from "@fluentui/react-icons";

export default function BlocksList() {
  const { providers } = useProjectContext();

  return (
      <div className="absolute flex flex-col gap-1 m-3 z-10 top-8" style={{marginTop: "40px"}}>
        <Accordion defaultOpenItems={"Built-in"} collapsible={true}>
          {providers.map((r) => {
            const p = r.provider;
            if (!p) {
                return null;
            }
            const resError = r.info && r.info.state === ProviderState.ERROR;
            return (
                <AccordionItem key={p.id} value={p.name} disabled={resError}>
                  <AccordionHeader icon={resError ? <PlugDisconnected20Regular /> : null}>
                    {p.name}
                  </AccordionHeader>
                  <AccordionPanel className={"overflow-y-auto"} style={{maxHeight: "40em"}}>
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

