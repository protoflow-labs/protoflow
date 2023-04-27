import { Divider, Field, Input } from "@fluentui/react-components";
import { useForm } from "react-hook-form";
import { Node } from "reactflow";
import { EditorActions, useUnselect } from "../EditorActions";

export type BucketData = {
  name: string;
  config: { bucket?: { path: string } };
};

export function BucketEditor(props: { node: Node<BucketData> }) {
  const onCancel = useUnselect();
  const { register, handleSubmit, watch } = useForm({
    values: {
      name: props.node.data.name,
      path: props.node.data.config.bucket?.path || "",
    },
  });

  const onSubmit = (data: any) => {
    props.node.data.name = data.name;

    if (!props.node.data.config.bucket) {
      props.node.data.config.bucket = {
        path: "",
      };
    }

    props.node.data.config.bucket.path = data.path;

    onCancel();
  };

  const values = watch();

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="flex flex-col gap-2 p-3">
        <Field label="Name" required>
          <Input value={values.name} {...register("name")} />
        </Field>
        <Field label="Path">
          <Input value={values.path} {...register("path")} />
        </Field>
        <Divider />
        <EditorActions />
      </div>
    </form>
  );
}
