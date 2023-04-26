import { InputData } from "@/components/blocks/InputBlock";
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
import { FieldDefinition, FieldType } from "../../../rpc/block_pb";
import { EditorActions, useUnselect } from "../EditorActions";
import { fieldTypeMap } from "../EditorPanel";

export function InputEditor({ node }: { node: Node<InputData> }) {
  const onCancel = useUnselect();
  const { watch, setValue, register, handleSubmit } = useForm({
    values: {
      name: node.data.name || "",
      fields: node.data.config.input?.fields || [],
    },
  });

  const onSubmit = (data: any) => {
    node.data.name = data.name;

    if (!node.data.config.input) {
      node.data.config.input = {
        fields: [],
      };
    }

    node.data.config.input.fields = data.fields;

    onCancel();
  };

  const values = watch();

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="flex flex-col gap-2 p-3">
        <Field label="Name" required>
          <Input value={values.name} {...register("name")} />
        </Field>
        <div className="flex flex-col">
          <Field label="Fields">
            {values.fields?.map((field, index) => (
              <div key={index} className="flex items-center gap-2 mb-2">
                <Input
                  value={values.fields[index].name}
                  {...register(`fields.${index}.name`)}
                />
                <Dropdown
                  id={"fieldType" + index}
                  value={fieldTypeMap[field.type || FieldType.STRING]}
                  onOptionSelect={(_, data) => {
                    setValue(`fields.${index}.type`, Number(data.optionValue));
                  }}
                >
                  <Option value={String(FieldType.STRING)}>String</Option>
                  <Option value={String(FieldType.INTEGER)}>Integer</Option>
                  <Option value={String(FieldType.BOOLEAN)}>Boolean</Option>
                </Dropdown>
                <Button
                  icon={<HiOutlineTrash className="h-4 w-4" />}
                  onClick={() => {
                    setValue(
                      "fields",
                      values.fields.filter((_, i) => i !== index)
                    );
                  }}
                />
              </div>
            ))}
          </Field>

          <Button
            icon={<HiPlus className="w-4 h-4" />}
            onClick={() => {
              setValue("fields", [
                ...(values.fields || []),
                new FieldDefinition({
                  name: "",
                  type: FieldType.STRING,
                }),
              ]);
            }}
          >
            Add Field
          </Button>
        </div>
        <Divider />
        <EditorActions />
      </div>
    </form>
  );
}
