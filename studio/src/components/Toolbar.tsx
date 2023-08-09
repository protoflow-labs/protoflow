import { generateService, projectService } from "@/lib/api";
import { checkIsApple } from "@/lib/checkIsApple";

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

export function Toolbar() {
  const isApple = checkIsApple();
  const { project, runWorkflow } = useProjectContext();
  const { save, props, setMode } = useEditorContext();
  const prevNodeLength = useRef(0);

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
  });


  const onSave = async () => {
    await save();
  }

  const onBuild = async () => {
    await onSave();

    generateService.generate({ projectId: project?.id });
    toast.success("Project built");
  };

  const onEditor = async () => {
    setMode("editor");
  }

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
          </MenuList>
        </MenuPopover>
      </Menu>

      <Button appearance="subtle" size="small" onClick={onEditor}>
       Editor
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
