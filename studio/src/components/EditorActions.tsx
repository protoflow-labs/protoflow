import { Button } from "@fluentui/react-components";
import { ReactFlowState, useStore } from "reactflow";
import {generateService} from "@/lib/api";
import {useProjectContext} from "@/providers/ProjectProvider";

const selectResetSelectedElements = (state: ReactFlowState) =>
  state.resetSelectedElements;
export const useUnselect = () => useStore(selectResetSelectedElements);

export function EditorActions() {
  const onCancel = useUnselect();
  const { project } = useProjectContext();
    const { activeNode } = useProjectContext();

  const updateType = () => {
    generateService.inferNodeType({ projectId: project?.id, nodeId: activeNode?.id });
  }

  const buildNode = () => {
    generateService.generateImplementation({ projectId: project?.id, nodeId: activeNode?.id });
  }

  return (
      <>
          <div className="flex items-center justify-between gap-2">
              <Button onClick={updateType}>Update Type</Button>
              <Button onClick={buildNode}>Implementation</Button>
          </div>
          <div className="flex items-center justify-between gap-2">
              <Button onClick={onCancel}>Cancel</Button>
              <Button appearance="primary" type="submit">
                  Save
              </Button>
          </div>
      </>
  );
}
