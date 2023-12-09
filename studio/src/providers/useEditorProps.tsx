import {
    applyEdgeChanges, applyNodeChanges,
    Connection,
    Edge,
    MarkerType,
    Node,
    OnEdgesChange,
    OnNodesChange,
    ReactFlowInstance
} from "reactflow";
import {useProjectContext} from "@/providers/ProjectProvider";
import {DragEventHandler, useCallback, useState} from "react";
import {v4 as uuid} from "uuid";
import {generateUUID} from "@/util/uuid";
import {DraggedNode} from "@/providers/EditorProvider";
import {Edge as ProtoEdge, Map, Node as ProtoNode, NodeDetails, Provides} from "@/rpc/graph_pb";
import {BaseNode} from "@/components/BaseNode";

const nodeTypes: Record<string, any> = {
    'node': BaseNode,
};

export const useEditorProps = (
    draggedNode: DraggedNode | undefined,
    setDraggedNode: (node: DraggedNode) => void,
    setSelectedNodes: (nodes: ProtoNode[]) => void,
    setSelectedEdges: (edges: ProtoEdge[]) => void,
    reactFlowInstance?: ReactFlowInstance,
) => {
    const {project, setNodeLookup, setEdgeLookup, edgeLookup, nodeLookup} = useProjectContext();

    const [nodes, setNodes] = useState<Node[]>(
        project?.graph?.nodes.map((n) => {
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

            // TODO breadchris what is going on in this code?
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
