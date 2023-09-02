import {createContext, DragEventHandler, ReactNode, useCallback, useContext, useEffect, useState,} from "react";
import {
    Connection,
    Edge,
    Node,
    OnEdgesChange,
    OnNodesChange,
    ReactFlowInstance, ReactFlowState, useStore,
} from "reactflow";
import {useProjectContext} from "./ProjectProvider";
import {Edge as ProtoEdge, Map, Node as ProtoNode, NodeDetails, Provides} from "@/rpc/graph_pb";
import {projectService} from "@/lib/api";
import { GetNodeInfoResponse } from "@/rpc/project_pb";
import {useEditorProps} from "@/providers/useEditorProps";

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
    setDraggedNode: (node: DraggedNode) => void;
    setInstance: (instance: ReactFlowInstance) => void;
    setMode: (mode: Mode) => void;
    save: () => void;

    selectedNodes: ProtoNode[];
    selectedEdges: ProtoEdge[];
    nodeInfo: GetNodeInfoResponse | undefined;
};

type Mode = "editor" | "run";

const EditorContext = createContext<EditorContextType>({} as any);

export const useEditorContext = () => useContext(EditorContext);
export const useEditorMode = () => useEditorContext().mode;
export const useActivelyEditing = () => {
    const {selectedNodes, selectedEdges} = useEditorContext();
    return selectedNodes.length > 0 || selectedEdges.length > 0;
}
export const useCurrentNode = () => {
    const {selectedNodes} = useEditorContext();
    return selectedNodes.length === 0 ? undefined : selectedNodes[0];
}
export const useUnselect = () => useStore((state: ReactFlowState) => state.resetSelectedElements);

export interface DraggedNode {
    provider: NodeDetails;
    node: ProtoNode;
}

export function EditorProvider({children}: { children: ReactNode }) {
    const {saveProject, nodeLookup, edgeLookup, project} = useProjectContext();

    const [draggedNode, setDraggedNode] = useState<DraggedNode | undefined>(undefined);
    const [instance, setInstance] = useState<ReactFlowInstance>();
    const [mode, setMode] = useState<Mode>("editor");

    const [selectedNodes, setSelectedNodes] = useState<ProtoNode[]>([]);
    const [selectedEdges, setSelectedEdges] = useState<ProtoEdge[]>([]);
    const [nodeInfo, setNodeInfo] = useState<GetNodeInfoResponse | undefined>(undefined);

    //useAutoLayout({direction: 'TB'})

    useEffect(() => {
        if (!project || selectedNodes.length !== 1) return;
        const selectedNode = selectedNodes[0];
        if (!selectedNode) {
            setNodeInfo(undefined);
            console.error(`selected node is undefined`)
            return;
        }

        (async () => {
            try {
                const info = await projectService.getNodeInfo({
                    nodeId: selectedNode.id,
                    projectId: project.id,
                });
                setNodeInfo(info);
            } catch (e) {
                // this is ok if we error, the node might not exist yet
                console.warn(e);
            }
        })();
    }, [project, selectedNodes]);

    const props = useEditorProps(
        draggedNode,
        setDraggedNode,
        setSelectedNodes,
        setSelectedEdges,
        instance,
    );

    const save = useCallback(async () => {
        return await saveProject(props.nodes, props.edges);
    }, [props.nodes, props.edges, saveProject]);

    useEffect(() => {
        void save();
    }, [nodeLookup, edgeLookup]);

    return (
        <EditorContext.Provider value={{
            props,
            mode,
            save,
            setMode,
            setInstance,
            setDraggedNode,
            selectedNodes,
            selectedEdges,
            nodeInfo,
        }}>
            {children}
        </EditorContext.Provider>
    );
}
