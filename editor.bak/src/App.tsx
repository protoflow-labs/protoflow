import {
  Button,
  Dropdown,
  FluentProvider,
  Option,
  webDarkTheme,
} from "@fluentui/react-components";
import { compile as hbs } from "handlebars";
import { DragEventHandler, useCallback, useState } from "react";
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
import { v4 as uuid } from "uuid";
import "./App.css";
import { EditorPanel } from "./components/EditorPanel";
import { EntityNode } from "./nodes/EntityNode";
import { FunctionNode } from "./nodes/FunctionNode";
import { InputNode } from "./nodes/InputNode";
import { ValidatorNode } from "./nodes/ValidatorNode";

import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from "@tanstack/react-query";
import BlocksList from "./components/BlocksList";
import DefaultEdge from "./edges/DefaultEdge";
import { projectService } from "./lib/api";
import { NodeResourceDependencies } from "./lib/resources";
import { BucketNode } from "./nodes/BucketNode";
import { QueryNode } from "./nodes/QueryNode";
import { QueueNode } from "./nodes/QueueNode";
import InputEntityEdgeTemplate from "./templates/InputEntityEdgeTemplate.hbs?raw";

const queryClient = new QueryClient();
const generateInputEntityEdgeTemplate = hbs(InputEntityEdgeTemplate);

const initialNodes: any = [];
const initialEdges: any = [];

const nodeTypes = {
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
  }, []);

  const onNodesDelete: OnNodesDelete = useCallback((nodes) => {}, []);

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

    setNodes((nds) => [...nds, newNode]);
  };

  const onDragOver: DragEventHandler = useCallback((e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
  }, []);

  const onNodeDragStop: NodeDragHandler = useCallback((e, node) => {}, []);

  return (
    <QueryClientProvider client={queryClient}>
      <FluentProvider theme={webDarkTheme}>
        <ReactFlowProvider>
          <div id="app">
            <div className="flex flex-col gap-4 absolute m-4 bg-white/10 p-4 rounded-md z-10">
              <Projects />

              <BlocksList />
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
