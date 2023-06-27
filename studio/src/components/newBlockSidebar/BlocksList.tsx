import {
  Accordion,
  AccordionHeader,
  AccordionItem,
  AccordionPanel,
  Button,
  Popover,
  Tooltip
} from "@fluentui/react-components";
import { useProjectContext } from "@/providers/ProjectProvider";
import {GRPC, Function, Input, Collection, Bucket} from "@/rpc/block_pb";
import {Node} from "@/rpc/graph_pb";
import {Resource} from "@/rpc/resource_pb";
import {NodeButton} from "@/components/newBlockSidebar/NodeButton";
import { ResourceState } from "@/rpc/project_pb";
import { PlugDisconnected20Regular } from "@fluentui/react-icons";

interface NodeBlock {
  type: string
  name: string
}

function resourceToNode(res: Resource, name: string) {
  const baseNode = new Node({
    name: name,
    resourceId: res.id,
  });

  switch (res.type.case) {
    case 'languageService':
      baseNode.config = {
        case: 'function',
        value: new Function({
          runtime: res.type.value.runtime,
        })
      }
      break;
    case 'docstore':
      baseNode.config = {
        case: 'collection',
        value: new Collection({})
      }
      break;
    case 'blobstore':
      baseNode.config = {
        case: 'bucket',
        value: new Bucket({})
      }
      break;
    default:
      return null;


  }
  return baseNode;
}

export default function BlocksList() {
  const { resources, deleteResource } = useProjectContext();
  const builtinBlocks: Node[] = [
    new Node({
      name: "Input",
      config: {
        case: "input",
        value: new Input({})
      }
    })
  ];

  return (
      <div className="absolute flex flex-col gap-1 m-3 z-10 top-8" style={{marginTop: "40px"}}>
        <Accordion defaultOpenItems={"Built-in"} collapsible={true}>
          <AccordionItem value="Built-in">
            <AccordionHeader>Built-in</AccordionHeader>
            <AccordionPanel>
              {builtinBlocks.map((node, i) => {
                return (
                    <NodeButton key={i} node={node} />
                );
              })}
            </AccordionPanel>
          </AccordionItem>
          {resources.map((r) => {
            if (!r.resource || !r.resource.type || !r.resource.type.case) {
              return null;
            }
            const resError = r.info && r.info.state === ResourceState.ERROR;
            const res = r.resource;
            const n = resourceToNode(res, "new");
            return (
                <AccordionItem key={res.id} value={res.name} disabled={resError}>
                  <AccordionHeader icon={resError ? <PlugDisconnected20Regular /> : null}>
                    <Tooltip content={r.info?.error || r.resource.type.case} relationship={"description"}>
                      <div>{res.name}</div>
                    </Tooltip>
                  </AccordionHeader>
                  <AccordionPanel className={"overflow-y-auto"} style={{maxHeight: "40em"}}>
                    {r.nodes.length === 0 && (
                        <div className="text-gray-400">No nodes</div>
                    )}
                    {n && (
                        <NodeButton node={n} />
                    )}
                    {r.nodes.map((node) => {
                      // TODO breadchris support more block types
                      return (
                          <NodeButton key={node.id} node={node} />
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

