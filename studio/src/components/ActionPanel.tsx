import { useEditorMode } from "@/providers/EditorProvider";
import { EditorPanel } from "./EditorPanel";
import RunPanel from "./RunPanel";
import React from "react";
import ChatPanel from "@/components/chat";

export function ActionPanel() {
  const editorMode = useEditorMode();

  switch (editorMode) {
    case "editor":
      return <EditorPanel />;
    case "run":
      return <RunPanel />;
  }
}
