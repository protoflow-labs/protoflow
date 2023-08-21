import {Badge, Button, Card, Divider, Textarea} from "@fluentui/react-components";
import {useProjectContext} from "@/providers/ProjectProvider";
import {ProtoViewer} from "@/components/ProtoViewer";
import {EditorActions} from "@/components/EditorActions";
import React, {useEffect, useState} from "react";
import {GRPCInputFormProps, ProtobufInputForm} from "@/components/ProtobufInputForm";
import { GRPCTypeInfo } from "@/rpc/project_pb";
import {useForm} from "react-hook-form";
import {Node as ProtoNode, Edge as ProtoEdge} from "@/rpc/graph_pb";
import {useEditorContext} from "@/providers/EditorProvider";
import {toast} from "react-hot-toast";

interface EdgeEditorProps {
  edge: ProtoEdge;
}

const ActiveEdgeEditor: React.FC<EdgeEditorProps> = ({edge}) => {
  const { projectTypes , setEdgeLookup} = useProjectContext();
  const [config, setConfig] = useState<string>(JSON.stringify(edge.toJson()?.valueOf() || {}, null, 2));
  const { register, control, handleSubmit, setValue} = useForm({
    values: {
      data: edge.toJson()?.valueOf() || {}
    },
  });
  if (!projectTypes || !projectTypes.edgeType) {
    return null;
  }

  const inputFormProps: GRPCInputFormProps = {
    grpcInfo: new GRPCTypeInfo({
      input: projectTypes.edgeType,
      output: projectTypes.edgeType,
      descLookup: projectTypes.descLookup,
      enumLookup: projectTypes.enumLookup,
      packageName: '',
    }),
    // some random key to separate data from the form
    baseFieldName: 'data',
    //@ts-ignore
    register,
    setValue,
    // TODO breadchris without this ignore, my computer wants to take flight https://github.com/react-hook-form/react-hook-form/issues/6679
    //@ts-ignore
    control,
    fieldPath: '',
  }

  const onSubmit = async (data: any) => {
    setEdgeLookup((lookup) => {
      const edge = ProtoEdge.fromJson(data.data);
      return {
        ...lookup,
        [edge.id]: edge,
      }
    })
    toast.success('Saved!');
  };

  const saveConfig = () => {
    setEdgeLookup((lookup) => {
      const edge = ProtoEdge.fromJson(JSON.parse(config));
      return {
        ...lookup,
        [edge.id]: edge,
      }
    })
    toast.success('Saved!');
  }

  return (
      <div className="flex flex-col gap-2 p-3">
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="flex flex-col gap-2 p-3">
            <ProtobufInputForm {...inputFormProps} />
          </div>
          <div className="flex items-center">
            <Button appearance="primary" type="submit">
              Save
            </Button>
          </div>
          <Divider/>
          <Textarea value={config} onChange={(e) => setConfig(e.target.value)} />
          <Button onClick={saveConfig}>Save</Button>
        </form>
      </div>
  );
}

interface NodeEditorProps {
  node: ProtoNode;
}

// TODO breadchris mostly an exact duplicate of ActiveEdgeEditor, should be refactored
const HackyNodeEditor: React.FC<NodeEditorProps> = ({node}) => {
  const { projectTypes , setNodeLookup} = useProjectContext();
  const [config, setConfig] = useState<string>(JSON.stringify(node.toJson()?.valueOf() || {}, null, 2));
  const { register, control, handleSubmit, setValue} = useForm({
    values: {
      data: node.toJson()?.valueOf() || {}
    },
  });
  if (!projectTypes || !projectTypes.edgeType) {
    return null;
  }

  const inputFormProps: GRPCInputFormProps = {
    grpcInfo: new GRPCTypeInfo({
      input: projectTypes.nodeType,
      output: projectTypes.nodeType,
      descLookup: projectTypes.descLookup,
      enumLookup: projectTypes.enumLookup,
      packageName: '',
    }),
    // some random key to separate data from the form
    baseFieldName: 'data',
    //@ts-ignore
    register,
    setValue,
    // TODO breadchris without this ignore, my computer wants to take flight https://github.com/react-hook-form/react-hook-form/issues/6679
    //@ts-ignore
    control,
    fieldPath: '',
  }

  const onSubmit = async (data: any) => {
    setNodeLookup((lookup) => {
      const node = ProtoNode.fromJson(data.data);
      return {
        ...lookup,
        [node.id]: node,
      }
    })
    toast.success('Saved!');
  };

  const saveConfig = () => {
    setNodeLookup((lookup) => {
      const node = ProtoNode.fromJson(JSON.parse(config));
      return {
        ...lookup,
        [node.id]: node,
      }
    })
    toast.success('Saved!');
  }

  return (
      <div className="flex flex-col gap-2 p-3">
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="flex flex-col gap-2 p-3">
            <ProtobufInputForm {...inputFormProps} />
          </div>
          <div className="flex items-center">
            <Button appearance="primary" type="submit">
              Save
            </Button>
          </div>
          <Divider/>
          <Textarea value={config} onChange={(e) => setConfig(e.target.value)} />
          <Button onClick={saveConfig}>Save</Button>
        </form>
      </div>
  );
}

const ActiveNodeEditor: React.FC<NodeEditorProps> = ({node}) => {
  return (
      <>
        <HackyNodeEditor node={node} />
        <Divider/>
        <EditorActions/>
        <ProtoViewer />
      </>
  );
}

export function EditorPanel() {
  const { selectedNodes, selectedEdges } = useEditorContext();

  if (selectedNodes.length > 0) {
    return <ActiveNodeEditor node={selectedNodes[0]} />;
  } else if (selectedEdges.length > 0) {
    return <ActiveEdgeEditor edge={selectedEdges[0]} />;
  }
  return null;
}
