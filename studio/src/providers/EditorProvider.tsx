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
import {getDataFromNode, getNodeDataKey, useProjectContext} from "./ProjectProvider";
import {Block} from "@/rpc/block_pb";

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

export const ReactFlowProtoflowKey = "application/reactflow";

export type ReactFlowProtoflowData = {
  type: string;
  name: string | undefined;
  config: Block['type'] | undefined;
  resourceIds: string[];
}

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
          resourceIds: n.resourceIds,
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

    setEdges((eds) => {
      return addEdge({ ...params, id: uuid() }, eds)
    });
  }, [nodes]);

  const onDragOver: DragEventHandler = useCallback((e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
  }, []);

  const onDrop: DragEventHandler<HTMLDivElement> = useCallback(
    (e) => {
      e.preventDefault();

      const data = e.dataTransfer.getData(ReactFlowProtoflowKey);
      const { type, name, config, resourceIds } = JSON.parse(data) as ReactFlowProtoflowData;

      const position = reactFlowInstance!.project({
        x: e.clientX,
        y: e.clientY,
      });

      // TODO breadchris protobuf type expects object in the form { grpc: { ... } }
      const configType = config && config.case ? {
        [config.case]: config.value
      } : {};

      const newNode = {
        id: uuid(),
        type,
        position,
        data: {
          name: name || "",
          config: configType,
          resourceIds,
        },
      };

      setNodes((nds) => [...nds, newNode]);
    },
    [reactFlowInstance]
  );

  const onEdgesChange: OnEdgesChange = useCallback((changes) => {
    console.log('edge change called')
    setEdges((eds) => applyEdgeChanges(changes, eds));
  }, []);

  const onNodesChange: OnNodesChange = useCallback((changes) => {
    console.log('node change called', changes)
    setNodes((nds) => applyNodeChanges(changes, nds));
    console.log('nodes are now ', nodes)
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
