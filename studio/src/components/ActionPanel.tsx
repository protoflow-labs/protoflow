import { useEditorMode } from "@/providers/EditorProvider";
import { EditorPanel } from "./EditorPanel";
import RunPanel from "./RunPanel";

export function ActionPanel() {
  const editrMode = useEditorMode();

  switch (editrMode) {
    case "editor":
      return <EditorPanel />;

    case "run":
      return <RunPanel />;
  }
}
