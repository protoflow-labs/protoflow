import { ReactFlowState, useStore } from "reactflow";

interface Selected {
    selectedNodes: Node[];
    selectedEdges: ReactFlowState["edges"];
}

const selectSelectedNodes = (state: ReactFlowState) => {
  return {
    selectedNodes: state.getNodes().filter((node) => node.selected),
    selectedEdges: state.edges.filter((edge) => edge.selected),
  }
}

export const useSelectedNodes = () => {
  const selected = useStore(selectSelectedNodes);

  return {
    selectedNodes: selected.selectedNodes,
    selectedEdges: selected.selectedEdges,
  };
};
