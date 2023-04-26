import { EntityData } from "@/components/blocks/EntityBlock";
import { Divider, Field, Input } from "@fluentui/react-components";
import { useForm } from "react-hook-form";
import { Node } from "reactflow";
import { EditorActions, useUnselect } from "../EditorActions";

export function EntityEditor(props: { node: Node<EntityData> }) {
  const onCancel = useUnselect();
  const { register, handleSubmit, watch } = useForm({
    values: {
      name: props.node.data.name || "",
      collection: props.node.data.config.entity?.collection || "",
    },
  });

  const onSubmit = (data: any) => {
    props.node.data.name = data.name;

    if (!props.node.data.config.entity) {
      props.node.data.config.entity = {
        collection: "",
      };
    }

    props.node.data.config.entity.collection = data.collection;

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
          <Input value={values.collection} {...register("collection")} />
        </Field>
        <Divider />
        <EditorActions />
      </div>
    </form>
  );
}
