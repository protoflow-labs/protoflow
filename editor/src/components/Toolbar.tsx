import { useSelectedNodes } from "@/hooks/useSelectedNodes";
import { generateService, projectService } from "@/lib/api";
import { checkIsApple } from "@/lib/checkIsApple";
import { getUpdatedProject } from "@/lib/project";
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
import { useHotkeys } from "react-hotkeys-hook";

export function Toolbar() {
  const isApple = checkIsApple();
  const { project } = useProjectContext();
  const { props } = useEditorContext();
  const { selectedNodes } = useSelectedNodes();

  useHotkeys(isApple ? "meta+s" : "ctrl+s", (e) => {
    e.preventDefault();
    e.stopPropagation();
    onSave();
  });

  useHotkeys(isApple ? "meta+b" : "ctrl+b", (e) => {
    e.preventDefault();
    e.stopPropagation();
    onBuild();
  });

  useHotkeys(isApple ? "meta+enter" : "ctrl+enter", (e) => {
    e.preventDefault();
    e.stopPropagation();
    onRun();
  });

  const onSave = async () => {
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

    await projectService.saveProject(updatedProject);
    toast.success("Project saved");
  };

  const onBuild = async () => {
    await onSave();

    generateService.generate({ projectId: project?.id });
    toast.success("Project built");
  };

  const onRun = async () => {
    if (selectedNodes.length !== 1) {
      toast.error("Please select a node to run");
      return;
    }

    const selectedNode = selectedNodes[0];
    const res = projectService.runWorklow({
      projectId: project?.id,
      nodeId: selectedNode.id,
    });

    console.log(res);
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
            <MenuItem secondaryContent="Ctrl+S" onClick={onSave}>
              Save
            </MenuItem>
            <MenuItem secondaryContent="Ctrl+Shift+E">Export</MenuItem>
          </MenuList>
        </MenuPopover>
      </Menu>

      <Button appearance="subtle" size="small">
        Edit
      </Button>
      <Button appearance="subtle" size="small">
        View
      </Button>
      <Menu>
        <MenuTrigger disableButtonEnhancement>
          <Button appearance="subtle" size="small">
            Build
          </Button>
        </MenuTrigger>

        <MenuPopover>
          <MenuList>
            <MenuItem secondaryContent="Ctrl+B" onClick={onBuild}>
              Build
            </MenuItem>
          </MenuList>
        </MenuPopover>
      </Menu>
      <Menu>
        <MenuTrigger disableButtonEnhancement>
          <Button appearance="subtle" size="small">
            Run
          </Button>
        </MenuTrigger>

        <MenuPopover>
          <MenuList>
            <MenuItem
              secondaryContent="Ctrl+R"
              disabled={selectedNodes.length !== 1}
              onClick={onRun}
            >
              Run Workflow
            </MenuItem>
            <MenuItem secondaryContent="Ctrl+B">Execute Block</MenuItem>
          </MenuList>
        </MenuPopover>
      </Menu>
    </div>
  );
}
