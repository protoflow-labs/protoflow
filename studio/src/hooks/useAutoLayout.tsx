import { useEffect } from 'react';
import { Node, Edge, Position, ReactFlowState, useStore, useReactFlow } from 'reactflow';
import { stratify, tree } from 'd3-hierarchy';

// the layout direction (T = top, R = right, B = bottom, L = left, TB = top to bottom, ...)
export type Direction = 'TB' | 'LR' | 'RL' | 'BT';

export type Options = {
    direction: Direction;
};

const positionMap: Record<string, Position> = {
    T: Position.Top,
    L: Position.Left,
    R: Position.Right,
    B: Position.Bottom,
};

const getPosition = (x: number, y: number, direction: Direction) => {
    switch (direction) {
        case 'LR':
            return { x: y, y: x };
        case 'RL':
            return { x: -y, y: -x };
        case 'BT':
            return { x: -x, y: -y };
        default:
            return { x, y };
    }
};

// initialize the tree layout (see https://observablehq.com/@d3/tree for examples)
const layout = tree<Node>()
    // the node size configures the spacing between the nodes ([width, height])
    .nodeSize([130, 120])
    // this is needed for creating equal space between all nodes
    .separation(() => 1);

const nodeCountSelector = (state: ReactFlowState) => state.nodeInternals.size;
const nodesInitializedSelector = (state: ReactFlowState) =>
    Array.from(state.nodeInternals.values()).every((node) => node.width && node.height);

function useAutoLayout(options: Options) {
    const { direction } = options;
    const nodeCount = useStore(nodeCountSelector);
    const nodesInitialized = useStore(nodesInitializedSelector);
    const { getNodes, getEdges, setNodes, setEdges, fitView } = useReactFlow();

    useEffect(() => {
        // only run the layout if there are nodes and they have been initialized with their dimensions
        if (!nodeCount || !nodesInitialized) {
            return;
        }

        const nodes: Node[] = getNodes();
        const edges: Edge[] = getEdges();

        // Find nodes that don't have any predecessor
        const nodesWithoutPredecessor = nodes.filter(
            node => !edges.some(edge => edge.target === node.id)
        );

        const FAKE_ROOT_ID = "FAKE_ROOT";

        // If we have any nodes without a predecessor, add a "fake root"
        if (nodesWithoutPredecessor.length > 0) {
            nodes.push({id: FAKE_ROOT_ID, type: "default", data: {}, position: {x: 0, y: 0}});
        }

        const hierarchy = stratify<Node>()
            .id((d) => d.id)
            .parentId((d: Node) => {
                // If this node doesn't have a predecessor, return the fake root ID
                if (nodesWithoutPredecessor.includes(d)) {
                    return FAKE_ROOT_ID;
                }
                // Otherwise, get the node's parent from the edges
                return edges.find((e: Edge) => e.target === d.id)?.source;
            })(nodes);

        // run the layout algorithm with the hierarchy data structure
        const root = layout(hierarchy);

        // set the React Flow nodes with the positions from the layout
        setNodes((nodes) =>
            nodes.map((node) => {
                // find the node in the hierarchy with the same id and get its coordinates
                const { x, y } = root.find((d) => d.id === node.id) || { x: node.position.x, y: node.position.y };

                return {
                    ...node,
                    sourcePosition: positionMap[direction[1]],
                    targetPosition: positionMap[direction[0]],
                    position: getPosition(x, y, direction),
                    style: { opacity: 1 },
                };
            })
        );

        setEdges((edges) => edges.map((edge) => ({ ...edge, style: { opacity: 1 } })));
    }, [nodeCount, nodesInitialized, getNodes, getEdges, setNodes, setEdges, fitView, direction]);
}

export default useAutoLayout;
