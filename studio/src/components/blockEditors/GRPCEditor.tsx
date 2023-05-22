import {
  Badge,
  Divider,
  Field,
  Input,
} from "@fluentui/react-components";
import {useForm} from "react-hook-form";
import {Node} from "reactflow";
import {GRPC} from "@/rpc/block_pb";
import {EditorActions, useUnselect} from "../EditorActions";
import {GRPCData} from "@/components/blocks/GRPCBlock";
import {getNodeDataKey, useProjectContext} from "@/providers/ProjectProvider";
import React, {useEffect, useState} from "react";
import {GetNodeInfoResponse} from "@/rpc/project_pb";
import SyntaxHighlighter from 'react-syntax-highlighter';


export function GRPCEditor({node}: { node: Node<GRPCData> }) {
  const onCancel = useUnselect();
  const {resources, loadNodeInfo} = useProjectContext();
  const {watch, setValue, register, handleSubmit, control} = useForm({
    values: {
      name: node.data.name || "",
      config: {
        package: node.data.config.grpc?.package || "",
        service: node.data.config.grpc?.service || "",
        method: node.data.config.grpc?.method || "",
      } as GRPC,
      input: node.data.config.grpc?.typeInfo?.input,
      data: JSON.parse(localStorage.getItem(getNodeDataKey(node)) || '{}'),
    },
  });
  const values = watch();
  const [nodeInfo, setNodeInfo] = useState<GetNodeInfoResponse | undefined>(undefined);

  // TODO breadchris this should be more general and not just for grpc
  useEffect(() => {
    (async () => {
      const res = await loadNodeInfo(node.id);
      setNodeInfo(res);
    })()
  }, [node.id]);

  const onSubmit = (data: any) => {
    node.data.name = data.name;

    if (!node.data.config.grpc) {
      node.data.config.grpc = {
        package: "",
        service: "",
        method: "",
        // TODO breadchris this is not a valid grpc descriptor proto, but OK for now?
        // @ts-ignore
        input: {},
        data: {},
      };
    }

    node.data.config.grpc.package = data.config.package;
    node.data.config.grpc.service = data.config.service;
    node.data.config.grpc.method = data.config.method;

    localStorage.setItem(getNodeDataKey(node), JSON.stringify(data.data));
    onCancel();
  };

  const fieldPath = node.data.config.grpc.package;
  //@ts-ignore
  const inputFormProps: GRPCInputFormProps = {
    grpcInfo: node.data.config.grpc.typeInfo,
    // some random key to separate data from the form
    baseFieldName: 'data',
    //@ts-ignore
    register,
    // TODO breadchris without this ignore, my computer wants to take flight https://github.com/react-hook-form/react-hook-form/issues/6679
    //@ts-ignore
    control,
    fieldPath,
  }
  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="flex flex-col gap-2 p-3">
        <Field label="Name" required>
          <Input value={values.name} {...register("name")} />
        </Field>
        <Field label="Package">
          <Input value={values.config.package} {...register("config.package")} />
        </Field>
        <Field label="Service" required>
          <Input value={values.config.service} {...register("config.service")} />
        </Field>
        <Field label="Method" required>
          <Input value={values.config.method} {...register("config.method")} />
        </Field>
        <Divider/>
        {resources && resources.filter((r) => node.data.resourceIds?.indexOf(r.id) >= 0).map((r) => (
          <Badge key={r.id}>{r.name}</Badge>
        ))}
        <Divider/>
        {nodeInfo && (
            <div className={"max-w-sm"}>
              <SyntaxHighlighter language="protobuf">
                {nodeInfo.methodProto}
              </SyntaxHighlighter>
            </div>
        )}
        <Divider/>
        <EditorActions/>
      </div>
    </form>
  );
}
