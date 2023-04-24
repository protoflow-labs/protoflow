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

export function Toolbar() {
  const { project } = useProjectContext();
  const {
    props: { edges, nodes },
  } = useEditorContext();

  const onExport = () => {
    if (!project) return;
    const updatedProject = {
      projectId: project.id,
      graph: {
        edges: edges.map((edge) => ({
          id: edge.id,
          source: edge.source,
          target: edge.target,
        })),
        nodes: nodes.map((node) => ({
          id: node.id,
          name: node.data.name,
          x: node.position.x,
          y: node.position.y,
        })),
      },
      // projectId: project?.getProjectId(),
    };

    for (const node of updatedProject.graph.nodes) {
      if (!node.name) {
        toast.error("Please name all nodes before exporting");
      }
    }
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
