import {Badge, Card, Divider} from "@fluentui/react-components";
import {useOnSelectionChange, useReactFlow} from "reactflow";
import NodeProvider from "@/providers/NodeProvider";
import {getNodeDataKey, useProjectContext} from "@/providers/ProjectProvider";
import {ActionBar} from "@/components/ActionBar";
import {ProtoViewer} from "@/components/blockEditors/ProtoViewer";
import {EditorActions} from "@/components/EditorActions";
import React, {useEffect, useState} from "react";
import {useSelectedNodes} from "@/hooks/useSelectedNodes";
import {GRPCInputFormProps, ProtobufInputForm} from "@/components/inputForms/ProtobufInputForm";
import { GRPCTypeInfo } from "@/rpc/project_pb";
import {useForm} from "react-hook-form";

export function EditorPanel() {
  const { projectTypes } = useProjectContext();
  const { resources, activeNode, activeEdge, setActiveNodeId, setActiveEdgeId } = useProjectContext();
  const {watch, setValue, register, handleSubmit, control} = useForm({
    values: {
      from: activeEdge ? activeEdge.from : '',
      to: activeEdge ? activeEdge.to : '',
    },
  });
  const values = watch();

  useOnSelectionChange({
    onChange: ({ nodes, edges }) => {
      if (nodes.length !== 1) {
        setActiveNodeId(null);
      } else {
        setActiveNodeId(nodes[0].id);
      }
      if (edges.length !== 1) {
        setActiveEdgeId(null);
      } else {
        setActiveEdgeId(edges[0].id);
      }
    },
  });

  if (activeNode) {
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

  if (!projectTypes || !projectTypes.edgeType) {
    return null;
  }

  const inputFormProps: GRPCInputFormProps = {
    grpcInfo: new GRPCTypeInfo({
      input: projectTypes.edgeType,
      output: projectTypes.edgeType,
      packageName: '',
    }),
    // some random key to separate data from the form
    baseFieldName: 'data',
    //@ts-ignore
    register,
    // TODO breadchris without this ignore, my computer wants to take flight https://github.com/react-hook-form/react-hook-form/issues/6679
    //@ts-ignore
    control,
    fieldPath: '',
  }

  if (activeEdge) {
    const a = projectTypes && projectTypes.edgeType
    return (
      <div className="absolute top-0 right-0 m-4 z-10 overflow-auto">
        <Card>
          <ProtobufInputForm {...inputFormProps} />
        </Card>
      </div>
    );
  }
  return null;
}
