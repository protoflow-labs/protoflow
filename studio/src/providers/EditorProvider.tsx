import { useBlockTypes } from "@/hooks/useBlockTypes";
import { configTypes } from "@/lib/configTypes";
import {
  createContext,
  DragEventHandler,
  ReactNode,
  useCallback,
  useContext,
  useState,
} from "react";
import {
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
  Connection,
  Edge,
  Node,
  OnEdgesChange,
  OnNodesChange,
  ReactFlowInstance,
} from "reactflow";
import { v4 as uuid } from "uuid";
import { useProjectContext } from "./ProjectProvider";

type EditorContextType = {
  mode: Mode;
  props: {
    edges: Edge[];
    nodes: Node[];
    nodeTypes: Record<string, any>;
    onConnect: (params: Connection) => void;
    onDragOver: DragEventHandler;
    onDrop: DragEventHandler<HTMLDivElement>;
    onEdgesChange: OnEdgesChange;
    onNodesChange: OnNodesChange;
  };
  setInstance: (instance: ReactFlowInstance) => void;
  setMode: (mode: Mode) => void;
};

type Mode = "editor" | "run";

const EditorContext = createContext<EditorContextType>({} as any);

export const useEditorContext = () => useContext(EditorContext);
export const useEditorMode = () => useEditorContext().mode;

export function EditorProvider({ children }: { children: ReactNode }) {
  const [instance, setInstance] = useState<ReactFlowInstance>();
  const [mode, setMode] = useState<Mode>("editor");
  const props = useEditorProps(instance);

  return (
    <EditorContext.Provider value={{ props, mode, setMode, setInstance }}>
      {children}
    </EditorContext.Provider>
  );
}

const useEditorProps = (reactFlowInstance?: ReactFlowInstance) => {
  const { project } = useProjectContext();

  const [nodes, setNodes] = useState<Node[]>(
    project?.graph?.nodes.map((n) => {
      const config = configTypes.find((c) => n.config?.case === c.name);

      return {
        id: n.id,
        data: {
          name: n.name,
          config: {
            [config!.name]:
              n.config?.value?.toJson() || n.config?.value || n.config || {},
          },
        },
        position: { x: n.x, y: n.y },
        type: `protoflow.${config?.name}`,
      };
    }) || []
  );

  const [edges, setEdges] = useState<Edge[]>(
    project?.graph?.edges?.map((e) => ({
      id: e.id,
      source: e.from,
      target: e.to,
    })) || []
  );

  const { nodeTypes } = useBlockTypes();

  const onConnect = useCallback((params: Connection) => {
    if (!params.source || !params.target) return;

    setEdges((eds) => addEdge({ ...params, id: uuid() }, eds));
  }, []);

  const onDragOver: DragEventHandler = useCallback((e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
  }, []);

  const onDrop: DragEventHandler<HTMLDivElement> = useCallback(
    (e) => {
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
        data: { name: "", config: {} },
      };

      setNodes((nds) => [...nds, newNode]);
    },
    [reactFlowInstance]
  );

  const onEdgesChange: OnEdgesChange = useCallback((changes) => {
    setEdges((eds) => applyEdgeChanges(changes, eds));
  }, []);

  const onNodesChange: OnNodesChange = useCallback((changes) => {
    setNodes((nds) => applyNodeChanges(changes, nds));
  }, []);

  return {
    edges,
    nodes,
    nodeTypes,
    onConnect,
    onDragOver,
    onDrop,
    onNodesChange,
    onEdgesChange,
  };
};