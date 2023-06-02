import { useSelectedNodes } from "@/hooks/useSelectedNodes";
import { generateService, projectService } from "@/lib/api";
import { checkIsApple } from "@/lib/checkIsApple";
import { saveProject } from "@/lib/project";
import { useEditorContext } from "@/providers/EditorProvider";
import { useProjectContext } from "@/providers/ProjectProvider";
import {
  Button,
  Menu,
  MenuItem,
  MenuList,
  MenuPopover,
  MenuTrigger,
  Dialog
} from "@fluentui/react-components";
import {useCallback, useEffect, useRef, useState} from "react";
import { toast } from "react-hot-toast";
import { useHotkeys } from "react-hotkeys-hook";
import {AddResourceDialog} from "@/components/AddResourceDialog";
import { useErrorBoundary } from "react-error-boundary";

export function Toolbar() {
  const isApple = checkIsApple();
  const { project, runWorkflow } = useProjectContext();
  const { props } = useEditorContext();
  const { selectedNodes } = useSelectedNodes();
  const prevNodeLength = useRef(0);
  const [addResourceOpen, setAddResourceOpen] = useState(false);

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

  useHotkeys("shift+enter", (e) => {
    e.preventDefault();
    e.stopPropagation();
    onRun();
  });

  useHotkeys("ctrl+shift+a", (e) => {
    e.preventDefault();
    e.stopPropagation();
    onAddResource();
  });

  const onSave = useCallback(async () => {
    if (!project) return;

    const updatedProject = saveProject({
      project,
      nodes: props.nodes,
      edges: props.edges,
    });

    for (const node of updatedProject.graph?.flowNodes || []) {
      if (!node.name) {
        toast.error("Please name all nodes before exporting");
        return;
      }
    }

    await projectService.saveProject(updatedProject);
    toast.success("Project saved");
  }, [project, props.flowNodes, props.edges]);

  const onBuild = async () => {
    await  onSave();

    generateService.generate({ projectId: project?.id });
    toast.success("Project built");
  };

  const onRun = async () => {
    if (selectedNodes.length !== 1) {
      toast.error("Please select a node to run");
      return;
    }

    const selectedNode = selectedNodes[0];
    await runWorkflow(selectedNode);
  };

  const onAddResource = async () => {
    setAddResourceOpen(true);
  }

  useEffect(() => {
    if (
      selectedNodes.length === 0 &&
      selectedNodes.length !== prevNodeLength.current
    ) {
      onSave();
    }

    prevNodeLength.current = selectedNodes.length;
  }, [selectedNodes, onSave]);

  return (
    <div className="px-1 py-1 absolute z-10 top-0">
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
            <MenuItem secondaryContent="Ctrl+Shift+E">
              Export
            </MenuItem>
            <MenuItem secondaryContent="Ctrl+Shift+A" onClick={onAddResource}>
              Add Resource
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
      <AddResourceDialog open={addResourceOpen} close={() => setAddResourceOpen(false)} />
    </div>
  );
}
