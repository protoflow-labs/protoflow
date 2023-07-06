import { useEditorMode } from "@/providers/EditorProvider";
import { EditorPanel } from "./EditorPanel";
import RunPanel from "./RunPanel";
import React from "react";
import ChatPanel from "@/components/chat";

export function ActionPanel() {
  const editrMode = useEditorMode();

  switch (editrMode) {
    case "chat":
      return <ChatPanel />;
    case "editor":
      return <EditorPanel />;
    case "run":
      return <RunPanel />;
  }
}
