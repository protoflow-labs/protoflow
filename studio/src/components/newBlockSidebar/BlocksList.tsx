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
import {GRPC, Function, Input, Collection, Bucket, Prompt, Config, Query} from "@/rpc/block_pb";
import {Node} from "@/rpc/graph_pb";
import {Resource} from "@/rpc/resource_pb";
import {NodeButton} from "@/components/newBlockSidebar/NodeButton";
import { ResourceState } from "@/rpc/project_pb";
import { PlugDisconnected20Regular } from "@fluentui/react-icons";
import {useState} from "react";
import {projectService} from "@/lib/api";

interface NodeBlock {
  type: string
  name: string
}

function defaultNodesForResource(res: Resource, name: string) {
  const newNode = (config: Node['config']) => {
    return new Node({
      name: name,
      resourceId: res.id,
      config: config,
    });
  }
  let nodes: Node[] = [];

  switch (res.type.case) {
    case 'languageService':
      nodes.push(newNode({
        case: 'function',
        value: new Function({
          runtime: res.type.value.runtime,
        })
      }));
      break;
    case 'docStore':
      nodes.push(newNode({
        case: 'collection',
        value: new Collection({})
      }));
      nodes.push(newNode({
        case: 'query',
        value: new Query({})
      }));
      break;
    case 'fileStore':
      nodes.push(newNode({
        case: 'bucket',
        value: new Bucket({})
      }));
      break;
    case 'reasoningEngine':
      nodes.push(newNode({
        case: 'prompt',
        value: new Prompt({})
      }));
      break;
    case 'configProvider':
      nodes.push(newNode({
        case: 'configuration',
        value: new Config({})
      }));
      break;
    default:
      return nodes;
  }
  return nodes;
}

export default function BlocksList() {
  const { resources, deleteResource, updateResource } = useProjectContext();
  const [connectResource, setConnectResource] = useState<string>('');
  const builtinBlocks: Node[] = [
    new Node({
      name: "Input",
      config: {
        case: "input",
        value: new Input({})
      }
    }),
  ];

  const setResourceToConnect = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setConnectResource(e.target.value);
  }

  const doConnectResource = async (resource: Resource) => {
    resource.dependencies.push(connectResource);
    await updateResource(resource);
  }

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
            const n = defaultNodesForResource(res, "new");
            return (
                <AccordionItem key={res.id} value={res.name} disabled={resError}>
                  <AccordionHeader icon={resError ? <PlugDisconnected20Regular /> : null}>
                    <Tooltip content={r.info?.error || r.resource.type.case} relationship={"description"}>
                      <div>{res.name}</div>
                    </Tooltip>
                  </AccordionHeader>
                  <AccordionPanel className={"overflow-y-auto"} style={{maxHeight: "40em"}}>
                    {n.length > 0 ? (
                      <>
                        {n.map((node, idx) => (<NodeButton key={idx} node={node} />))}
                      </>
                    ) : (
                      <>
                        {r.nodes.length === 0 && (
                            <div className="text-gray-400">No nodes</div>
                        )}
                      </>
                    )}
                    {r.nodes.map((node) => {
                      // TODO breadchris support more block types
                      return (
                          <NodeButton key={node.id} node={node} />
                      );
                    })}
                    <Select onChange={setResourceToConnect}>
                      {resources.map((r) => {
                        if (!r.resource) {
                          return null;
                        }
                        return (
                            <option key={r.resource.id} value={r.resource.id}>{r.resource?.name}</option>
                        )
                      })}
                    </Select>
                    <Button size="small" className="w-full" appearance={'outline'} onClick={() => doConnectResource(res)}>Connect</Button>
                    <Button size="small" className="w-full" appearance={'outline'} onClick={() => deleteResource(res.id)}>Delete</Button>
                  </AccordionPanel>
                </AccordionItem>
            );
          })}
        </Accordion>
      </div>
  );
}

