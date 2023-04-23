import {
  Button,
  Dropdown,
  FluentProvider,
  Option,
  webDarkTheme,
} from "@fluentui/react-components";
import { compile as hbs } from "handlebars";
import {
  DragEventHandler,
  ReactNode,
  useCallback,
  useEffect,
  useState,
} from "react";
import ReactFlow, {
  Background,
  Connection,
  Edge,
  Node,
  NodeDragHandler,
  OnEdgesChange,
  OnNodesChange,
  OnNodesDelete,
  ReactFlowInstance,
  ReactFlowProvider,
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
} from "reactflow";

import "reactflow/dist/style.css";
import "./App.css";
import { EditorPanel } from "./components/EditorPanel";
import { EntityNode } from "./nodes/EntityNode";
import { FunctionNode } from "./nodes/FunctionNode";
import { InputNode } from "./nodes/InputNode";
import { ValidatorNode } from "./nodes/ValidatorNode";
import { v4 as uuid } from "uuid";

import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from "@tanstack/react-query";
import DefaultEdge from "./edges/DefaultEdge";
import { projectService } from "./lib/api";
import { NodeResourceDependencies } from "./lib/resources";
import { BucketNode } from "./nodes/BucketNode";
import { EndpointyNode } from "./nodes/EndpointNode";
import { QueryNode } from "./nodes/QueryNode";
import { QueueNode } from "./nodes/QueueNode";
import InputEntityEdgeTemplate from "./templates/InputEntityEdgeTemplate.hbs?raw";

const queryClient = new QueryClient();
const generateInputEntityEdgeTemplate = hbs(InputEntityEdgeTemplate);

const initialNodes: any = [];
const initialEdges: any = [];

const nodeTypes = {
  endpoint: EndpointyNode,
  entity: EntityNode,
  function: FunctionNode,
  validation: ValidatorNode,
  message: InputNode,
  query: QueryNode,
  queue: QueueNode,
  bucket: BucketNode,
};

const edgeTypes = {
  edge: DefaultEdge,
};

function App() {
  const [reactFlowInstance, setReactFlowInstance] =
    useState<ReactFlowInstance | null>(null);
  const [nodes, setNodes] = useState<Node<any>[]>(initialNodes);
  const [edges, setEdges] = useState<Edge<any>[]>(initialEdges);

  const onNodesChange: OnNodesChange = useCallback((changes) => {
    setNodes((nds) => applyNodeChanges(changes, nds));
  }, []);

  const onEdgesChange: OnEdgesChange = useCallback((changes) => {
    setEdges((eds) => applyEdgeChanges(changes, eds));
    changes.forEach((change => {
      if (change.type === 'remove') {
        projectService.removeEdge({
          edgeId: change.id
        })
      }
    }));
  }, []);

  const onNodesDelete: OnNodesDelete = useCallback((nodes) => {
    console.log(nodes);
    nodes.forEach(node => projectService.removeBlock({
      blockId: node.id,
    }));
  }, []);

  const onConnect = useCallback((params: Connection) => {
    if (!params.source || !params.target) return;

    setEdges((eds) =>
      addEdge(
        {
          ...params,
          type: "edge",
          data: { async: false },
        },
        eds
      )
    );
    projectService.addEdge({
      edge: {
        id: uuid(),
        source: params.source,
        target: params.target,
      }
    })
  }, []);

  const onDrop: DragEventHandler<HTMLDivElement> = (e) => {
    e.preventDefault();

    const type = e.dataTransfer.getData("application/reactflow");
    const position = reactFlowInstance!.project({
      x: e.clientX,
      y: e.clientY,
    });

    const newNode = {
      id: uuid(),
      type,
      position,
      data: { label: `${type} node` },
    };

    projectService.addBlock({
      block: {
        id: newNode.id,
        type: newNode.type,
        x: newNode.position.x,
        y: newNode.position.y,
      },
    });

    setNodes((nds) => [...nds, newNode]);
  };

  const onDragOver: DragEventHandler = useCallback((e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
  }, []);

  const onNodeDragStop: NodeDragHandler = useCallback((e, node) => {
    projectService.updateBlock({
      block: {
        id: node.id,
        x: node.position.x,
        y: node.position.y,
        name: node.data.name,
        type: node.type,
      },
    });
  }, []);

  useEffect(() => {
    (async function () {
      projectService.getBlocks({}).then(({ blocks }) => {
        const nodes = blocks.map((block) => {
          return {
            id: block.id,
            type: block.type,
            position: { x: block.x, y: block.y },
            data: { name: block.name },
          };
        });
        setNodes(nodes);
      });

      projectService.getEdges({}).then(({ edges }) => {
        setEdges(edges.map((edge) => {
          return {
            id: edge.id,
            source: edge.source,
            target: edge.target,
          };
        }));
      });
    })();
  }, []);

  return (
    <QueryClientProvider client={queryClient}>
      <FluentProvider theme={webDarkTheme}>
        <ReactFlowProvider>
          <div id="app">
            <div className="flex flex-col gap-4 absolute m-4 bg-white/10 p-4 rounded-md z-10">
              <Projects />

              <div className="flex flex-col gap-1">
                <NodeButton nodeType="endpoint">Endpoint</NodeButton>
                <NodeButton nodeType="entity">Entity</NodeButton>
                <NodeButton nodeType="message">Input</NodeButton>
                <NodeButton nodeType="function">Function</NodeButton>
                <NodeButton nodeType="query">Query</NodeButton>
                <NodeButton nodeType="queue">Queue</NodeButton>
                <NodeButton nodeType="bucket">Bucket</NodeButton>
              </div>
            </div>

            <ReactFlow
              edges={edges}
              edgeTypes={edgeTypes}
              nodes={nodes}
              nodeTypes={nodeTypes}
              onConnect={onConnect}
              onDragOver={onDragOver}
              onNodesDelete={onNodesDelete}
              onDrop={onDrop}
              onEdgesChange={onEdgesChange}
              onInit={(ref: any) => setReactFlowInstance(ref)}
              onNodesChange={onNodesChange}
              onNodeDragStop={onNodeDragStop}
              proOptions={{ hideAttribution: true }}
              fitView
            >
              <Background />
            </ReactFlow>

            <Button
              size="small"
              className="absolute top-4 right-4"
              onClick={() => {
                const resources = new Set();
                for (const node of nodes) {
                  if (!node.type) continue;
                  const dependencies = NodeResourceDependencies[node.type];
                  if (!dependencies) continue;

                  for (const dependency of dependencies) {
                    resources.add(dependency);
                  }
                }

                const data = JSON.stringify(
                  { nodes, edges, resources: Array.from(resources) },
                  null,
                  2
                );
                console.log(data);

                // const dataStr =
                //   "data:text/json;charset=utf-8," + encodeURIComponent(data);
                // const link = document.createElement("a");
                // link.setAttribute("href", dataStr);
                // link.setAttribute("download", "protoflow-project.json");
                // document.body.appendChild(link); // required for firefox
                // link.click();
                // link.remove();
              }}
            >
              Export
            </Button>
            <EditorPanel />
          </div>
        </ReactFlowProvider>
      </FluentProvider>
    </QueryClientProvider>
  );
}

export default App;

function NodeButton(props: { children: ReactNode; nodeType: string }) {
  return (
    <div
      draggable
      onDragStart={(e) => {
        e.dataTransfer.setData("application/reactflow", props.nodeType);
        e.dataTransfer.effectAllowed = "move";
      }}
    >
      <Button>{props.children}</Button>
    </div>
  );
}

function Projects() {
  const { data } = useQuery({
    queryKey: ["projects"],
    queryFn: projectService.getProjects,
  });

  if (!data) return null;

  return (
    <Dropdown>
      {data.projects.map((project) => {
        return <Option key={project.name}>{project.name}</Option>;
      })}
    </Dropdown>
  );
}
