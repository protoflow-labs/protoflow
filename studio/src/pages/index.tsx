import { ActionBar } from "@/components/ActionBar";
import { ActionPanel } from "@/components/ActionPanel";
import BlocksList from "@/components/newBlockSidebar/BlocksList";
import { useEditorContext } from "@/providers/EditorProvider";
import { useProjectContext } from "@/providers/ProjectProvider";
import { Background, ReactFlow } from "reactflow";

import "reactflow/dist/style.css";

export default function Home() {
  const { project } = useProjectContext();
  const { props, setInstance } = useEditorContext();

  if (!project) {
    return null;
  }

  return (
    <>
      <main className="flex">
        <div className="flex flex-1">
          <BlocksList />
          <ActionPanel />
          <ReactFlow
            onInit={setInstance}
            proOptions={{ hideAttribution: true }}
            {...props}
            fitView
          >
            <Background />
          </ReactFlow>
          <ActionBar />
        </div>
      </main>
    </>
  );
}
