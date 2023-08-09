import {Button} from "@fluentui/react-components";
import {ReactFlowState, useStore} from "reactflow";
import {generateService} from "@/lib/api";
import {useProjectContext} from "@/providers/ProjectProvider";
import {useActivelyEditing, useCurrentNode, useEditorContext} from "@/providers/EditorProvider";

const selectResetSelectedElements = (state: ReactFlowState) =>
    state.resetSelectedElements;
export const useUnselect = () => useStore(selectResetSelectedElements);

export function EditorActions() {
    const {project} = useProjectContext();
    const activeNode = useCurrentNode();

    const updateType = () => {
        generateService.inferNodeType({projectId: project?.id, nodeId: activeNode?.id});
    }

    const buildNode = () => {
        generateService.generateImplementation({projectId: project?.id, nodeId: activeNode?.id});
    }

    return (
        <>
            <div className="flex items-center justify-between gap-2">
                <Button onClick={updateType}>Update Type</Button>
                <Button onClick={buildNode}>Implementation</Button>
            </div>
        </>
    );
}
