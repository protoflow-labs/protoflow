import { EntityData } from "@/blocks/EntityBlock";
import { Divider, Field, Input } from "@fluentui/react-components";
import { useForm } from "react-hook-form";
import { Node } from "reactflow";
import { EditorActions, useUnselect } from "../EditorActions";

export function EntityEditor(props: { node: Node<EntityData> }) {
  const onCancel = useUnselect();
  const { register, handleSubmit, watch } = useForm({
    values: {
      name: props.node.data.name || "",
      table: props.node.data.table || "",
    },
  });

  const onSubmit = (data: any) => {
    props.node.data.name = data.name;
    props.node.data.table = data.table;
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
          <Input value={values.table} {...register("table")} />
        </Field>
        <Divider />
        <EditorActions />
      </div>
    </form>
  );
}
