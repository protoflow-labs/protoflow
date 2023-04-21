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

const generateInputEntityEdgeTemplate = hbs(InputEntityEdgeTemplate);

const initialNodes: any = [
  // {
  //   id: "1",
  //   type: "entity",
  //   data: { table: "users", name: "Users" },
  //   position: { x: 400, y: 125 },
  // },
  // {
  //   id: "2",
  //   type: "input",
  //   data: { name: "CreateUserInput", fields: [] },
  //   position: { x: 400, y: 25 },
  // },
  // {
  //   id: "3",
  //   type: "default",
  //   data: { label: "SendWelcomeEmail" },
  //   position: { x: 250, y: 250 },
  // },
  // {
  //   id: "4",
  //   type: "default",
  //   data: { label: "TrackUserCreated" },
  //   position: { x: 400, y: 250 },
  // },
  // {
  //   id: "5",
  //   type: "default",
  //   data: { label: "SendSlackMessage" },
  //   position: { x: 550, y: 250 },
  // },
];

const initialEdges: any = [
  // { id: "e1-2", source: "2", target: "1", label:'Auto Implementation: insertRecordFromInput' },
  // { id: "e2-3", source: "1", target: "3", animated: true, label: "on INSERT" },
  // { id: "e2-4", source: "1", target: "4", animated: true, label: "on INSERT"  },
  // { id: "e2-5", source: "1", target: "5", animated: true, label: "on INSERT"   },
];

function App() {
  const nodeTypes = useMemo(
    () => ({
      entity: EntityNode,
      function: FunctionNode,
      validation: ValidatorNode,
      message: InputNode,
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
            <NodeButton nodeType="entity">Entity</NodeButton>
            <NodeButton nodeType="message">Input</NodeButton>
            <NodeButton nodeType="function">Function</NodeButton>
          </aside>
          <Button
            size="small"
            className="absolute top-4 right-4"
            onClick={() => {
              for (const edge of edges) {
                const sourceNode: Node<InputData> | undefined = nodes.find(
                  (node) => node.id === edge.source
                );

                const targetNode: Node<EntityData> | undefined = nodes.find(
                  (node) => node.id === edge.target
                );

                if (
                  sourceNode?.type === "message" &&
                  targetNode?.type === "entity"
                ) {
                  const template = generateInputEntityEdgeTemplate({
                    host: "docstore",
                    port: "27017",
                    database: "protoflow",
                    collection: targetNode.data.name,
                    name: sourceNode.data.name + "Impl",
                  });

                  console.log(template)
                }
              }
              // const data = JSON.stringify({ nodes, edges }, null, 2);

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
