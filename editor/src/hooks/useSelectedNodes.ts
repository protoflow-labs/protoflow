import {ReactFlowState, useStore} from "reactflow";

const selectSelectedNodes = (state: ReactFlowState) => state.getNodes().filter((node) => node.selected);

export const useSelectedNodes = () => {
  return useStore(selectSelectedNodes);
}
