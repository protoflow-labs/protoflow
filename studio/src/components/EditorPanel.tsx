import {Badge, Card, Divider} from "@fluentui/react-components";
import {useOnSelectionChange} from "reactflow";
import NodeProvider from "@/providers/NodeProvider";
import {useProjectContext} from "@/providers/ProjectProvider";
import {ActionBar} from "@/components/ActionBar";
import {ProtoViewer} from "@/components/blockEditors/ProtoViewer";
import {EditorActions} from "@/components/EditorActions";
import React, {useEffect, useState} from "react";

export function EditorPanel() {
  const { resources, activeNode, setActiveNodeId } = useProjectContext();

  useOnSelectionChange({
    onChange: ({ nodes }) => {
      if (nodes.length !== 1) {
        setActiveNodeId(null);
        return;
      }
      setActiveNodeId(nodes[0].id);
    },
  });

  if (!activeNode) {
    return null;
  }

  const getResourceBadge = () => {
    if (!resources) {
      return null;
    }
    const res = resources.find((r) => {
      return r.resource && activeNode.resourceId === r.resource.id
    })
    if (!res || !res.resource) {
      return null;
    }
    return <Badge key={res.resource.id}>{res.resource.name}</Badge>
  }

  return (
    <NodeProvider nodeId={activeNode.id}>
      <div className="absolute top-0 right-0 m-4 z-10 overflow-auto">
        <Card>
          <ProtoViewer />
          <Divider/>
          {getResourceBadge()}
          <Divider/>
          <EditorActions/>
        </Card>
      </div>
      <ActionBar node={activeNode} />
    </NodeProvider>
  );
}
