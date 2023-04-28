import { InputData } from "@/components/blocks/InputBlock";
import {
  Button,
  Checkbox,
  Divider,
  Dropdown,
  Field,
  Input,
  Option,
} from "@fluentui/react-components";
import { useState } from "react";
import { UseFormSetValue, useForm } from "react-hook-form";
import { HiOutlineTrash, HiPlus } from "react-icons/hi2";
import { TbFileText } from "react-icons/tb";
import { Node } from "reactflow";
import {
  FieldDefinition,
  FieldType
} from "../../../rpc/block_pb";
import { EditorActions, useUnselect } from "../EditorActions";

const fieldTypeToDisplay = {
  [FieldType.STRING]: "String",
  [FieldType.INTEGER]: "Integer",
  [FieldType.BOOLEAN]: "Boolean",
};

export const stringToFieldType = {
  "STRING": FieldType.STRING,
  "INTEGER": FieldType.INTEGER,
  "BOOLEAN": FieldType.BOOLEAN,
}

type Form = {
  name: string;
  config: {
    fields: FieldDefinition[]
    sampleData: {
      [key: string]: string | number | boolean;
    }
  }
}

export function InputEditor({ node }: { node: Node<InputData> }) {
  const [showSampleDataForm, setShowSampleDataForm] = useState(false);
  const onCancel = useUnselect();
  const sampleDataStorageKey = `${node.data.name}-sampleData`;
  const { watch, setValue, register, handleSubmit } = useForm<Form>({
    values: {
      name: node.data.name || "",
      config: {
        fields: node.data.config.input?.fields?.map(value => ({
          ...value,
          // TODO - Really really hacky, but connect seems to serialize enums differently
          // than the typescript expects
          type: ((stringToFieldType as any)[value.type as any] || value.type || 0)
        }  as FieldDefinition)) || [],
        sampleData: JSON.parse(localStorage.getItem(sampleDataStorageKey) || '{}')
      },
    },
  });

  const onSubmit = (data: any) => {
    node.data.name = data.name;

    if (!node.data.config.input) {
      node.data.config.input = {
        fields: [],
      };
    }

    node.data.config.input.fields = data.config.fields;
    localStorage.setItem(sampleDataStorageKey, JSON.stringify(data.config.sampleData));

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
            {values.config.fields?.map((field, index) => (
              <div key={index} className="flex items-center gap-2 mb-2">
                <Input
                  value={values.config.fields[index].name}
                  {...register(`config.fields.${index}.name`)}
                />
                <Dropdown
                  id={"fieldType" + index}
                  value={fieldTypeToDisplay[field.type || FieldType.STRING]}
                  onOptionSelect={(_, data) => {
                    setValue(
                      `config.fields.${index}.type`,
                      Number(data.optionValue)
                    );
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
                      "config.fields",
                      values.config.fields.filter((_, i) => i !== index)
                    );

                    const newSampleData = { ...values.config.sampleData };
                    delete newSampleData[field.name];
                    setValue('config.sampleData', newSampleData)
                  }}
                />
              </div>
            ))}
          </Field>

          <Button
            icon={<HiPlus className="w-4 h-4" />}
            onClick={() => {
              setValue("config.fields", [
                ...(values.config.fields || []),
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
        <TbFileText onClick={() => setShowSampleDataForm(!showSampleDataForm)}/>
        { showSampleDataForm && (
          <div>
            { values.config.fields.map((field, index) => (
              <div key={`${field.name}-${index}`} className="flex items-center gap-2 mb-2">
                <Field label={field.name}>
                  <SampleDataField field={field} values={values} setValue={setValue}/>
                </Field>
              </div>
            )) }
          </div>)}
        <Divider />
        <EditorActions />
      </div>
    </form>
  );
}

function SampleDataField({ field, values, setValue }: {
  field: FieldDefinition,
  values: Form,
  setValue: UseFormSetValue<Form>
}) {
  switch (field.type) {
    case FieldType.STRING:
      return <Input
          value={String(values.config.sampleData[field.name] || '')}
          onChange={(_, data) => setValue(`config.sampleData.${field.name}`, data.value)}
        />;
    case FieldType.INTEGER:
      return <Input
        type="number"
        value={String(values.config.sampleData[field.name])}
        onChange={(_, data) => setValue(`config.sampleData.${field.name}`, Number(data.value))}
      />;
    case FieldType.BOOLEAN:
      return <Checkbox
          checked={Boolean(values.config.sampleData[field.name])}
          onChange={(_, data) => setValue(`config.sampleData.${field.name}`, data.checked)}
        />;
  }
}
