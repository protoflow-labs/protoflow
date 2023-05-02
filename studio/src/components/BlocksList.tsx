import { Button } from "@fluentui/react-components";
import { ReactNode } from "react";

export default function BlocksList() {
  return (
    <div className="absolute flex flex-col gap-1 m-3 z-10 top-8">
      <NodeButton nodeType="protoflow.input">Input</NodeButton>
      <NodeButton nodeType="protoflow.collection">Collection</NodeButton>
      <NodeButton nodeType="protoflow.function">Function</NodeButton>
      <NodeButton nodeType="protoflow.query">Query</NodeButton>
      <NodeButton nodeType="protoflow.queue">Queue</NodeButton>
      <NodeButton nodeType="protoflow.bucket">Bucket</NodeButton>
      <NodeButton nodeType="protoflow.email">Email</NodeButton>
      <NodeButton nodeType="protoflow.rest">REST</NodeButton>
    </div>
  );
}

function NodeButton(props: { children: ReactNode; nodeType: string }) {
  return (
    <div
      draggable
      onDragStart={(e) => {
        e.dataTransfer.setData("application/reactflow", props.nodeType);
        e.dataTransfer.effectAllowed = "move";
      }}
    >
      <Button size="small" className="w-full">
        {props.children}
      </Button>
    </div>
  );
}