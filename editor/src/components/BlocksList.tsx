import { Button } from "@fluentui/react-components";
import { ReactNode } from "react";

export default function BlocksList() {
  return (
    <div className="absolute flex flex-col gap-1 m-3 z-10">
      <NodeButton nodeType="message">Input</NodeButton>
      <NodeButton nodeType="entity">Entity</NodeButton>
      <NodeButton nodeType="function">Function</NodeButton>
      <NodeButton nodeType="query">Query</NodeButton>
      <NodeButton nodeType="queue">Queue</NodeButton>
      <NodeButton nodeType="bucket">Bucket</NodeButton>
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
        console.log({ nodeType: props.nodeType });
      }}
    >
      <Button size="small" className="w-full">
        {props.children}
      </Button>
    </div>
  );
}
