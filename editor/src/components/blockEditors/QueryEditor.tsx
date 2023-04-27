import { Divider, Field, Input } from "@fluentui/react-components";
import { useForm } from "react-hook-form";
import { Node } from "reactflow";
import { Query } from "../../../rpc/block_pb";
import { EditorActions, useUnselect } from "../EditorActions";

export type QueryData = {
  name: string;
  config: { query?: Partial<Query> };
};

export function QueryEditor(props: { node: Node<QueryData> }) {
  const onCancel = useUnselect();
  const { register, handleSubmit, watch } = useForm({
    values: {
      name: props.node.data.name || "",
      config: {
        ...props.node.data.config.query,
      } as Query,
    },
  });

  const onSubmit = (data: any) => {
    props.node.data.name = data.name;

    if (!props.node.data.config.query) {
      props.node.data.config.query = {
        collection: "",
      };
    }

    props.node.data.config.query.collection = data.config.collection;

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
          <Input
            value={values.config.collection}
            {...register("config.collection")}
          />
        </Field>
        <Divider />
        <EditorActions />
      </div>
    </form>
  );
}
