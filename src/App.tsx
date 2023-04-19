import {
  DragEvent,
  DragEventHandler,
  MouseEvent,
  MouseEventHandler,
  useCallback,
  useMemo,
  useState,
} from "react";
import ReactFlow, {
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
  Background,
  Connection,
  Edge,
  Node,
  OnEdgesChange,
  OnNodesChange,
  ReactFlowInstance,
  ReactFlowProvider,
} from "reactflow";

import "reactflow/dist/style.css";
import "./App.css";
import { EntityNode } from "./nodes/EntityNode";
import { FunctionNode } from "./nodes/FunctionNode";
import { ValidatorNode } from "./nodes/ValidatorNode";
import { InputNode } from "./nodes/InputNode";

const initialNodes = [
  {
    id: "1",
    type: "entity",
    data: { table: "users", name: "Users" },
    position: { x: 400, y: 125 },
  },

  {
    id: "2",
    type: "input",
    data: { name: "CreateUserInput", fields: [] },
    position: { x: 400, y: 25 },
  },
  {
    id: "3",
    type: "default",
    data: { label: "SendWelcomeEmail" },
    position: { x: 250, y: 250 },
  },
  {
    id: "4",
    type: "default",
    data: { label: "TrackUserCreated" },
    position: { x: 400, y: 250 },
  },
  {
    id: "5",
    type: "default",
    data: { label: "SendSlackMessage" },
    position: { x: 550, y: 250 },
  },
];

const initialEdges = [
  { id: "e1-2", source: "2", target: "1" },
  { id: "e2-3", source: "1", target: "3", animated: true },
  { id: "e2-4", source: "1", target: "4", animated: true },
  { id: "e2-5", source: "1", target: "5", animated: true },
];

function App() {
  const nodeTypes = useMemo(
    () => ({
      entity: EntityNode,
      function: FunctionNode,
      validation: ValidatorNode,
      input: InputNode,
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
    console.log(params);
    setEdges((eds) => addEdge(params, eds));
  }, []);

  const onDrop: DragEventHandler<HTMLDivElement> = (event) => {
    event.preventDefault();

    const type = event.dataTransfer.getData("application/reactflow");
    const position = reactFlowInstance!.project({
      x: event.clientX,
      y: event.clientY,
    });

    const newNode = {
      id: `dndnode_${nodes.length}`,
      type,
      position,
      data: { label: `${type} node` },
    };

    setNodes((nds) => [...nds, newNode]);
  };

  const onDragOver: DragEventHandler = useCallback((event) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = "move";
  }, []);

  const onDragStart = (event: DragEvent<HTMLDivElement>, nodeType: string) => {
    event.dataTransfer.setData("application/reactflow", nodeType);
    event.dataTransfer.effectAllowed = "move";
  };

  return (
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
        <aside
          style={{
            position: "absolute",
            top: 0,
            left: 0,
            margin: 16,
            padding: 24,
            background: `rgba(255 ,255 ,255 ,0.1)`,
          }}
        >
          <div className="description">
            You can drag these nodes to the pane on the right.
          </div>
          <div
            className="dndnode entity"
            onDragStart={(event) => onDragStart(event, "entity")}
            draggable
          >
            Entity Node
          </div>
          <div
            className="dndnode"
            onDragStart={(event) => onDragStart(event, "input")}
            draggable
          >
            Input Node
          </div>
          <div
            className="dndnode"
            onDragStart={(event) => onDragStart(event, "validation")}
            draggable
          >
            Validator Node
          </div>
          <div
            className="dndnode function"
            onDragStart={(event) => onDragStart(event, "function")}
            draggable
          >
            Function Node
          </div>
        </aside>
      </div>
    </ReactFlowProvider>
  );
}

export default App;
