import {
  Button,
  Divider,
  Dropdown,
  Field,
  Input,
  Option,
} from "@fluentui/react-components";
import { useForm } from "react-hook-form";
import { HiOutlineTrash, HiPlus } from "react-icons/hi2";
import { Node } from "reactflow";
import {GRPC, REST} from "@/rpc/block_pb";
import { EditorActions, useUnselect } from "../EditorActions";
import { RESTData } from "@/components/blocks/RESTBlock";
import {GRPCData} from "@/components/blocks/GRPCBlock";

export function GRPCEditor({ node }: { node: Node<GRPCData> }) {
  const onCancel = useUnselect();
  const { watch, setValue, register, handleSubmit } = useForm({
    values: {
      name: node.data.name || "",
      config: {
        package: node.data.config.grpc?.package || "",
        service: node.data.config.grpc?.service || "",
        method: node.data.config.grpc?.method || "",
      } as GRPC,
    },
  });

  const onSubmit = (data: any) => {
    node.data.name = data.name;

    if (!node.data.config.grpc) {
      node.data.config.grpc = {
        package: "",
        service: "",
        method: "",
      };
    }

    node.data.config.grpc.package = data.config.package;
    node.data.config.grpc.service = data.config.service;
    node.data.config.grpc.method = data.config.method;

    onCancel();
  };

  const values = watch();

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
        <Divider />
        <EditorActions />
      </div>
    </form>
  );
}
