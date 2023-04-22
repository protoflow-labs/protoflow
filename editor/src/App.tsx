import {
  Button,
  FluentProvider,
  webDarkTheme,
} from "@fluentui/react-components";
import {
  DragEvent,
  DragEventHandler,
  ReactNode,
  useCallback,
  useMemo,
  useState,
} from "react";
import { compile as hbs } from "handlebars";
import ReactFlow, {
  Background,
  Connection,
  Edge,
  Node,
  OnEdgesChange,
  OnNodesChange,
  ReactFlowInstance,
  ReactFlowProvider,
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
} from "reactflow";

import "reactflow/dist/style.css";
import "./App.css";
import { EditorPanel } from "./components/EditorPanel";
import { EntityData, EntityNode } from "./nodes/EntityNode";
import { FunctionNode } from "./nodes/FunctionNode";
import { InputData, InputNode } from "./nodes/InputNode";
import { ValidatorNode } from "./nodes/ValidatorNode";

import InputEntityEdgeTemplate from "./templates/InputEntityEdgeTemplate.hbs?raw";
import { QueryNode } from "./nodes/QueryNode";
import { QueueNode } from "./nodes/QueueNode";
import { BucketNode } from "./nodes/BucketNode";
import { NodeResourceDependencies } from "./lib/resources";
import { EndpointyNode } from "./nodes/EndpointNode";

const generateInputEntityEdgeTemplate = hbs(InputEntityEdgeTemplate);

const initialNodes: any = [];
const initialEdges: any = [];

function App() {
  const nodeTypes = useMemo(
    () => ({
      endpoint: EndpointyNode,
      entity: EntityNode,
      function: FunctionNode,
      validation: ValidatorNode,
      message: InputNode,
      query: QueryNode,
      queue: QueueNode,
      bucket: BucketNode,
    }),
    []
  );

  const [reactFlowInstance, setReactFlowInstance] =
    useState<ReactFlowInstance | null>(null);
  const [nodes, setNodes] = useState<Node<any>[]>(initialNodes);
  const [edges, setEdges] = useState<Edge<any>[]>(initialEdges);

  const onNodesChange: OnNodesChange = useCallback(
    (changes) => setNodes((nds) => applyNodeChanges(changes, nds)),
    []
  );

  const onEdgesChange: OnEdgesChange = useCallback(
    (changes) => setEdges((eds) => applyEdgeChanges(changes, eds)),
    []
  );

  const onConnect = useCallback((params: Connection) => {
    setEdges((eds) => addEdge(params, eds));
  }, []);

  const onDrop: DragEventHandler<HTMLDivElement> = (e) => {
    e.preventDefault();

    const type = e.dataTransfer.getData("application/reactflow");

    const position = reactFlowInstance!.project({
      x: e.clientX,
      y: e.clientY,
    });

    const newNode = {
      id: `dndnode_${nodes.length}`,
      type,
      position,
      data: { label: `${type} node` },
    };

    setNodes((nds) => [...nds, newNode]);
  };

  const onDragStart = (event: DragEvent<HTMLDivElement>, nodeType: string) => {
    event.dataTransfer.setData("application/reactflow", nodeType);
    event.dataTransfer.effectAllowed = "move";
  };

  const onDragOver: DragEventHandler = useCallback((e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
  }, []);

  return (
    <FluentProvider theme={webDarkTheme}>
      <ReactFlowProvider>
        <div id="app">
          <ReactFlow
            nodeTypes={nodeTypes}
            nodes={nodes}
            onNodesChange={onNodesChange}
            edges={edges}
            onEdgesChange={onEdgesChange}
            onInit={(ref: any) => setReactFlowInstance(ref)}
            onDrop={onDrop}
            onDragOver={onDragOver}
            onConnect={onConnect}
            onChange={console.log}
            fitView
            proOptions={{ hideAttribution: true }}
          >
            <Background />
          </ReactFlow>
          <aside className="absolute top-0 left-0 m-4 p-2 bg-white/10 text-white rounded flex flex-col gap-2 items-start justify-start">
            <NodeButton nodeType="endpoint">Endpoint</NodeButton>
            <NodeButton nodeType="entity">Entity</NodeButton>
            <NodeButton nodeType="message">Input</NodeButton>
            <NodeButton nodeType="function">Function</NodeButton>
            <NodeButton nodeType="query">Query</NodeButton>
            <NodeButton nodeType="queue">Queue</NodeButton>
            <NodeButton nodeType="bucket">Bucket</NodeButton>
          </aside>
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
