import { useEditorMode } from "@/providers/EditorProvider";
import { EditorPanel } from "./EditorPanel";
import RunPanel from "./RunPanel";
import React from "react";
import ChatPanel from "@/components/chat";
import {Card, SelectTabData, SelectTabEvent, Tab, TabList, TabValue} from "@fluentui/react-components";
import {Stack} from "@fluentui/react";

export function ActionPanel() {
  const editorMode = useEditorMode();
    const [activeTab, setActiveTab] = React.useState<TabValue>('chat');
    const tabSelect = (event: SelectTabEvent, data: SelectTabData) => {
        setActiveTab(data.value);
    };

  switch (editorMode) {
    case "editor":
      return (
          <div className="absolute top-0 right-0 m-4 z-10 overflow-auto h-full">
              <Card>
                  <TabList onTabSelect={tabSelect} vertical={true}>
                      <Tab value='chat'>Chat</Tab>
                      <Tab value='workflow'>Workflow</Tab>
                  </TabList>
                  <Stack>
                      {activeTab === 'chat' ? <ChatPanel /> : <EditorPanel />}
                  </Stack>
              </Card>
          </div>
      );
    case "run":
      return <RunPanel />;
  }
}
