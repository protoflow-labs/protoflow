import {configTypes} from "@/lib/configTypes";
import {createContext, DragEventHandler, ReactNode, useCallback, useContext, useEffect, useState,} from "react";
import {
    applyEdgeChanges,
    applyNodeChanges,
    Connection,
    Edge,
    MarkerType,
    Node,
    OnEdgesChange,
    OnNodesChange,
    ReactFlowInstance,
} from "reactflow";
import {v4 as uuid} from "uuid";
import {useProjectContext} from "./ProjectProvider";
import {Edge as ProtoEdge, Map, Node as ProtoNode, Provides} from "@/rpc/graph_pb";
import {generateUUID} from "@/util/uuid";
import {StandardBlock} from "@/components/blocks/StandardBlock";
import {projectService} from "@/lib/api";
import { GetNodeInfoResponse } from "@/rpc/project_pb";
import useAutoLayout from "@/hooks/useAutoLayout";

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

export const ReactFlowProtoflowKey = "application/reactflow";

export type ReactFlowProtoflowData = {
    node: ProtoNode
}

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

const nodeTypes: Record<string, any> = {
    'node': StandardBlock,
};

export interface DraggedNode {
    provider: ProtoNode;
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

// todo: we want to make sure the incoming node has a distinction from the sidebar node vs the server type node
// export type SidebarNode = Exclude<ProtoNode, { id: string }>

const useEditorProps = (
    draggedNode: DraggedNode | undefined,
    setDraggedNode: (node: DraggedNode) => void,
    setSelectedNodes: (nodes: ProtoNode[]) => void,
    setSelectedEdges: (edges: ProtoEdge[]) => void,
    reactFlowInstance?: ReactFlowInstance,
) => {
    const {project, setNodeLookup, setEdgeLookup, edgeLookup, nodeLookup} = useProjectContext();

    const [nodes, setNodes] = useState<Node[]>(
        project?.graph?.nodes.map((n) => {
            const config = configTypes.find((c) => n.type?.case === c.name);

            return {
                id: n.id,
                data: {
                    node: n,
                },
                position: {x: n.x, y: n.y},
                type: 'node',
            };
        }) || []
    );

    const [edges, setEdges] = useState<Edge[]>(
        project?.graph?.edges?.map((e) => ({
            id: e.id,
            source: e.from,
            target: e.to,
            label: e.type?.case,
            markerEnd: {
                type: MarkerType.ArrowClosed
            },
        })) || []
    );

    const onConnect = useCallback((params: Connection) => {
        if (!params.source || !params.target) return;

        const newEdgeType: ProtoEdge = new ProtoEdge({
            id: uuid(),
            from: params.source,
            to: params.target,
            type: {
                case: 'map',
                value: new Map()
            }
        });
        const newEdge = {
            id: newEdgeType.id,
            source: newEdgeType.from,
            target: newEdgeType.to,
            label: newEdgeType.type.case,
        }

        setEdgeLookup((lookup) => {
            return {
                ...lookup,
                [newEdgeType.id]: newEdgeType,
            }
        })

        setEdges((eds) => [...eds, newEdge]);
    }, [nodes]);

    const onDragOver: DragEventHandler = useCallback((e) => {
        e.preventDefault();
        e.dataTransfer.dropEffect = "move";
    }, []);

    const onDrop: DragEventHandler<HTMLDivElement> = useCallback(
        (e) => {
            if (!draggedNode) {
                return;
            }
            const position = reactFlowInstance!.project({x: e.clientX, y: e.clientY});

            const newNode = {
                id: generateUUID(),
                type: 'node',
                position,
                data: {}
            };
            draggedNode.node.id = newNode.id;

            const newEdge = {
                id: uuid(),
                source: draggedNode.provider.id,
                target: draggedNode.node.id,
                label: "provides",
            }
            const newEdgeType: ProtoEdge = new ProtoEdge({
                id: newEdge.id,
                from: newEdge.source,
                to: newEdge.target,
                type: {
                    case: 'provides',
                    value: new Provides()
                }
            });

            setNodeLookup((lookup) => {
                return {
                    ...lookup,
                    [newNode.id]: draggedNode.node
                }
            })
            setEdgeLookup((lookup) => {
                return {
                    ...lookup,
                    [newEdge.id]: newEdgeType,
                }
            })
            setNodes((nds) => [...nds, newNode]);
            setEdges((eds) => [...eds, newEdge]);
        },
        [reactFlowInstance, draggedNode]
    );

    const onEdgesChange: OnEdgesChange = useCallback((changes) => {
        let newSelectedEdges: ProtoEdge[] = [];
        let clearSelectedEdges = false;
        changes.forEach((change) => {
            if (change.type === "select") {
                if (change.selected) {
                    const edge = edgeLookup[change.id];
                    newSelectedEdges.push(edge);
                } else {
                    clearSelectedEdges = true;
                }
            }
        });
        if (newSelectedEdges.length > 0) {
            setSelectedEdges(newSelectedEdges);
        } else {
            if (clearSelectedEdges) {
                setSelectedEdges([]);
            }
        }
        setEdges((eds) => applyEdgeChanges(changes, eds));
    }, [edgeLookup]);

    const onNodesChange: OnNodesChange = useCallback((changes) => {
        let newSelectedNodes: ProtoNode[] = [];
        let clearSelectedNodes = false;
        changes.forEach((change) => {
            if (change.type === "select") {
                if (change.selected) {
                    const node = nodeLookup[change.id];
                    newSelectedNodes.push(node);
                } else {
                    clearSelectedNodes = true;
                }
            }
        });
        if (newSelectedNodes.length > 0) {
            setSelectedNodes(newSelectedNodes);
        } else {
            if (clearSelectedNodes) {
                setSelectedNodes([]);
            }
        }
        setNodes((nds) => applyNodeChanges(changes, nds));
    }, [nodeLookup]);

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
