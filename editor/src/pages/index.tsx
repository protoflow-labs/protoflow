import { ActionBar } from "@/components/ActionBar";
import { ActionPanel } from "@/components/ActionPanel";
import BlocksList from "@/components/BlocksList";
import { useEditorContext } from "@/providers/EditorProvider";
import { useProjectContext } from "@/providers/ProjectProvider";
import Head from "next/head";
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
      <Head>
        <title>Protoflow</title>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
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
