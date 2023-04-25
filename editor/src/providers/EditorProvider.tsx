import { useNodeTypes } from "@/hooks/useNodeTypes";
import { configTypes } from "@/lib/configTypes";
import {
  DragEventHandler,
  ReactNode,
  createContext,
  useCallback,
  useContext,
  useState,
} from "react";
import {
  Connection,
  Edge,
  Node,
  OnEdgesChange,
  OnNodesChange,
  ReactFlowInstance,
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
} from "reactflow";
import { v4 as uuid } from "uuid";
import { useProjectContext } from "./ProjectProvider";

type EditorContextType = {
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
};

const EditorContext = createContext<EditorContextType>({} as any);

export const useEditorContext = () => useContext(EditorContext);

export function EditorProvider({ children }: { children: ReactNode }) {
  const [instance, setInstance] = useState<ReactFlowInstance>();
  const props = useEditorProps(instance);

  return (
    <EditorContext.Provider value={{ props, setInstance }}>
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
  const { nodeTypes } = useNodeTypes();

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
