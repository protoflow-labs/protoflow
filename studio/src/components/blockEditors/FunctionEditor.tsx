import { FunctionData } from "@/components/blocks/FunctionBlock";
import { Divider, Field, Input, Select } from "@fluentui/react-components";
import { useForm } from "react-hook-form";
import { Node } from "reactflow";
import { EditorActions, useUnselect } from "../EditorActions";

export function FunctionEditor({ node }: { node: Node<FunctionData> }) {
  const onCancel = useUnselect();
  const { watch, register, handleSubmit } = useForm({
    values: {
      name: node.data.name || "",
      runtime: node.data.config.function?.runtime || "node",
    },
  });

  const onSubmit = (data: any) => {
    node.data.name = data.name;

    if (!node.data.config.function) {
      node.data.config.function = {
        runtime: "",
      };
    }

    node.data.config.function.runtime = data.runtime;

    onCancel();
  };

  const values = watch();
  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="flex flex-col gap-2 p-3">
        <div className="flex flex-col">
          <Field label="Name" required>
            <Input value={values.name || ""} {...register("name")} />
          </Field>
        </div>
        <div className="flex flex-col">
          <Field label="Language" required>
            <Select value={values.runtime} {...register("runtime")}>
              <option value="node">Node</option>
              <option value="go">Go</option>
              <option value="python">Python</option>
            </Select>
          </Field>
        </div>
        <Divider />
        <EditorActions />
      </div>
    </form>
  );
}
