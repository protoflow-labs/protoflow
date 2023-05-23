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
import {useProjectContext} from "@/providers/ProjectProvider";
import React from "react";
import {ProtoViewer} from "@/components/blockEditors/common/ProtoViewer";


export function GRPCEditor({node}: { node: Node<GRPCData> }) {
  const onCancel = useUnselect();
  const {resources} = useProjectContext();
  const {watch, setValue, register, handleSubmit, control} = useForm({
    values: {
      name: node.data.name || "",
      config: {
        package: node.data.config.grpc?.package || "",
        service: node.data.config.grpc?.service || "",
        method: node.data.config.grpc?.method || "",
      } as GRPC,
    },
  });
  const values = watch();

  const onSubmit = (data: any) => {
    node.data.name = data.name;

    if (!node.data.config.grpc) {
      node.data.config.grpc = new GRPC({
        package: "",
        service: "",
        method: "",
      });
    }

    node.data.config.grpc.package = data.config.package;
    node.data.config.grpc.service = data.config.service;
    node.data.config.grpc.method = data.config.method;
    onCancel();
  };

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
        <ProtoViewer />
        <Divider/>
        <EditorActions/>
      </div>
    </form>
  );
}
