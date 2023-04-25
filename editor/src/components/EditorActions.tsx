import { Button } from "@fluentui/react-components";
import { ReactFlowState, useStore } from "reactflow";

const selectResetSelectedElements = (state: ReactFlowState) =>
  state.resetSelectedElements;
export const useUnselect = () => useStore(selectResetSelectedElements);

export function EditorActions() {
  const onCancel = useUnselect();

  return (
    <div className="flex items-center gap-2">
      <Button onClick={onCancel}>Cancel</Button>
      <Button appearance="primary" type="submit">
        Save
      </Button>
    </div>
  );
}
