import { Divider, Field, Input } from "@fluentui/react-components";
import { useForm } from "react-hook-form";
import { Node } from "reactflow";
import { Collection } from "../../../rpc/block_pb";
import { EditorActions, useUnselect } from "../EditorActions";

export type CollectionData = {
  name: string;
  config: { collection?: Partial<Collection> };
};

export function CollectionEditor(props: { node: Node<CollectionData> }) {
  const onCancel = useUnselect();
  const { register, handleSubmit, watch } = useForm({
    values: {
      name: props.node.data.name || "",
      config: {
        ...props.node.data.config.collection,
      } as Collection,
    },
  });

  const onSubmit = (data: any) => {
    props.node.data.name = data.name;

    if (!props.node.data.config.collection) {
      props.node.data.config.collection = {
        name: "",
      };
    }

    props.node.data.config.collection.name = data.config.name;

    onCancel();
  };

  const values = watch();

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="flex flex-col gap-2 p-3">
        <Field label="Name" required>
          <Input value={values.name} {...register("name")} />
        </Field>
        <Field label="Table" required>
          <Input value={values.config.name} {...register("config.name")} />
        </Field>
        <Divider />
        <EditorActions />
      </div>
    </form>
  );
}
