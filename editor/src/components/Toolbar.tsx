import { projectService } from "@/lib/api";
import { useEditorContext } from "@/providers/EditorProvider";
import { useProjectContext } from "@/providers/ProjectProvider";
import {
  Button,
  Menu,
  MenuItem,
  MenuList,
  MenuPopover,
  MenuTrigger,
} from "@fluentui/react-components";
import { toast } from "react-hot-toast";
import { getUpdatedProject } from "../lib/project";

export function Toolbar() {
  const { project } = useProjectContext();
  const { props } = useEditorContext();

  const onExport = () => {
    if (!project) return;

    const updatedProject = getUpdatedProject({
      project,
      nodes: props.nodes,
      edges: props.edges,
    });

    for (const node of updatedProject.graph?.nodes || []) {
      if (!node.name) {
        toast.error("Please name all nodes before exporting");
        return;
      }
    }

    projectService.saveProject(updatedProject);
  };

  return (
    <div className="px-1 py-1">
      <Menu>
        <MenuTrigger disableButtonEnhancement>
          <Button appearance="subtle" size="small">
            File
          </Button>
        </MenuTrigger>

        <MenuPopover>
          <MenuList>
            <MenuItem secondaryContent="Ctrl+S">Save</MenuItem>
            <MenuItem secondaryContent="Ctrl+Shift+E" onClick={onExport}>
              Export
            </MenuItem>
          </MenuList>
        </MenuPopover>
      </Menu>

      <Button appearance="subtle" size="small">
        Edit
      </Button>
      <Button appearance="subtle" size="small">
        View
      </Button>
    </div>
  );
}
