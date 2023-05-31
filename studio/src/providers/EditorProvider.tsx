import { useBlockTypes } from "@/hooks/useBlockTypes";
import { configTypes } from "@/lib/configTypes";
import {
  createContext,
  DragEventHandler,
  ReactNode,
  useCallback,
  useContext, useEffect,
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
import {getDataFromNode, getNodeDataKey, useProjectContext} from "./ProjectProvider";
import {Node as ProtoNode} from "@/rpc/graph_pb";
import {Simulate} from "react-dom/test-utils";
import drag = Simulate.drag;

type EditorContextType = {
  mode: Mode;
  save: () => Promise<void>;
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
  setDraggedNode: (node: ProtoNode) => void;
  setInstance: (instance: ReactFlowInstance) => void;
  setMode: (mode: Mode) => void;
};

export const ReactFlowProtoflowKey = "application/reactflow";

export type ReactFlowProtoflowData = {
  node: ProtoNode
}

type Mode = "editor" | "run";

const EditorContext = createContext<EditorContextType>({} as any);

export const useEditorContext = () => useContext(EditorContext);
export const useEditorMode = () => useEditorContext().mode;

export function EditorProvider({ children }: { children: ReactNode }) {
  const { saveProject, nodeLookup } = useProjectContext();
  const [draggedNode, setDraggedNode] = useState<ProtoNode | undefined>(undefined);
  const [instance, setInstance] = useState<ReactFlowInstance>();
  const [mode, setMode] = useState<Mode>("editor");
  const props = useEditorProps(
    draggedNode,
    setDraggedNode,
    instance,
  );

  const save = useCallback(async () => {
    return await saveProject(props.nodes, props.edges);
  }, [props]);

  useEffect(() => {
    void save();
  }, [nodeLookup]);

  return (
    <EditorContext.Provider value={{
      props,
      mode,
      setMode,
      setInstance,
      save,
      setDraggedNode,
    }}>
      {children}
    </EditorContext.Provider>
  );
}

const nodeToType = (node: ProtoNode) => {
  return `protoflow.${node.config.case}`;
}

const useEditorProps = (draggedNode: ProtoNode | undefined, setDraggedNode: (node: ProtoNode) => void, reactFlowInstance?: ReactFlowInstance) => {
  const { project, saveProject, setNodeLookup } = useProjectContext();

  const [nodes, setNodes] = useState<Node[]>(
    project?.graph?.nodes.map((n) => {
      const config = configTypes.find((c) => n.config?.case === c.name);

      return {
        id: n.id,
        data: {
          node: n,
        },
        position: { x: n.x, y: n.y },
        type: nodeToType(n),
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

    setEdges((eds) => {
      return addEdge({ ...params, id: uuid() }, eds)
    });
  }, [nodes]);

  const onDragOver: DragEventHandler = useCallback((e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
  }, []);

  const onDrop: DragEventHandler<HTMLDivElement>  = useCallback(
    (e) => {
      if (!draggedNode) {
        return;
      }
      const position = reactFlowInstance!.project({x: e.clientX, y: e.clientY});

      const newNode = {
        id: draggedNode.id,
        type: nodeToType(draggedNode),
        position,
        data: {}
      };
      setNodeLookup((lookup) => {
        return {
          ...lookup,
          [draggedNode.id]: draggedNode
        }
      })
      setNodes((nds) => [...nds, newNode]);
    },
    [reactFlowInstance, draggedNode]
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
