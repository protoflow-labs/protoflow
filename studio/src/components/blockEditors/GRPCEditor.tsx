import {
  Accordion,
  AccordionHeader,
  AccordionItem,
  AccordionPanel, Badge,
  Button,
  Divider,
  Field,
  Input, Select,
} from "@fluentui/react-components";
import {Control, FieldValues, useFieldArray, useForm, UseFormRegister, useWatch} from "react-hook-form";
import {Node} from "reactflow";
import {GRPC} from "@/rpc/block_pb";
import {EditorActions, useUnselect} from "../EditorActions";
import {GRPCData} from "@/components/blocks/GRPCBlock";
import {
  DescriptorProto,
  EnumDescriptorProto,
  FieldDescriptorProto,
  FieldDescriptorProto_Label
} from "@bufbuild/protobuf";
import {FC, useState} from "react";
import {getDataFromNode, getNodeDataKey, useProjectContext} from "@/providers/ProjectProvider";
import {useProjectResources} from "@/hooks/useProjectResources";

type GrpcFormField = {
  type: 'field'
  name: string
  field: FieldDescriptorProto
}

type GrpcFormOneof = {
  type: 'oneof'
  name: string
  fields: FieldDescriptorProto[]
}

type GrpcFormFieldType = GrpcFormField | GrpcFormOneof

const getFieldName = (baseFieldName: string | undefined, field: FieldDescriptorProto, idx?: number): string => {
  if (!baseFieldName) {
    return field.name || '';
  }
  if (idx !== undefined) {
    return `${baseFieldName}.${idx}.${field.name}`;
  }
  return `${baseFieldName}.${field.name}`;
}

interface InputFormContentsProps extends GRPCInputFormProps {
  field: FieldDescriptorProto
  index?: number
}

const InputFormContents: FC<InputFormContentsProps> = (props) => {
  const {
    baseFieldName,
    field,
    descLookup,
    enumLookup,
    register,
    control,
    index,
    fieldPath,
  } = props;
  const [visibleTypes, setVisibleTypes] = useState<string[]>([]);

  const fieldFormName = getFieldName(baseFieldName, field, index);
  const fieldValue = useWatch({
    control,
    name: fieldFormName,
    defaultValue: '',
  });

  if (field.typeName) {
    // field.typeName == .name.othername, remove the leading dot
    const typeName = field.typeName.substring(1);
    const fieldType = descLookup[typeName];
    if (fieldType) {
      return (
        <div key={field.number}>
          {visibleTypes.includes(typeName) ? (
            <>
              <GrpcInputForm
                {...props}
                baseFieldName={fieldFormName}
                fieldPath={`${fieldPath}.${field.name}`}
              />
              <Button onClick={() => {
                setVisibleTypes(visibleTypes.filter((t) => t !== typeName))
              }}>Close</Button>
            </>
          ) : (
            <Button onClick={() => setVisibleTypes([...visibleTypes, typeName])}>{typeName}</Button>
          )}
        </div>
      )
    }
  }
  // TODO breadchris this should be checking for an actual type, not a string, field.type is a string, not a number
  // @ts-ignore
  if (field.type === "TYPE_ENUM") {
    if (!field.typeName) {
      throw new Error("Enum field has no type name");
    }
    const enumType = enumLookup[`${fieldPath}.${field.name}`] || [];
    if (!enumType) {
      throw new Error(`Enum type ${fieldPath}.${field.name} not found in ${Object.keys(enumLookup)}`);
    }

    return (
      <>
        <label htmlFor={field.name}>{field.name}</label>
        <Select id={field.name}>
          {enumType.value.map((e) => (
            <option key={e.name} value={e.name}>{e.name}</option>
          ))}
        </Select>
      </>
    )
  }
  return (
    <Field key={field.number} label={field.name} required>
      <Input value={fieldValue} {...register(fieldFormName)} />
    </Field>
  )
}

type GRPCFormValues = {
  name: string
  config: {
    package: string
    service: string
    method: string
  },
  fields: GRPC
}

interface GRPCInputFormProps {
  desc: DescriptorProto | undefined
  descLookup: { [key: string]: DescriptorProto }
  enumLookup: { [key: string]: EnumDescriptorProto }
  register: UseFormRegister<any>
  control: Control
  fieldPath: string
  baseFieldName?: string
}

const GrpcInputForm: FC<GRPCInputFormProps> = (props) => {
  const {
    desc,
    control,
    baseFieldName,
    fieldPath
  } = props;

  if (!desc) {
    return null;
  }

  const formatField = (field: FieldDescriptorProto) => {
    if (!field.name) {
      console.error('Field has no name', field)
      return null;
    }

    const {fields: formFields, append, prepend, remove, swap, move, insert} = useFieldArray({
      control,
      name: baseFieldName || 'input',
    });

    const inputProps: InputFormContentsProps = {
      ...props,
      field,
      fieldPath: `${fieldPath}.${desc.name}`
    }

    // TODO breadchris for some reason FieldDescriptorProto_Label.REPEATED is a number and field.label is a string
    // @ts-ignore
    if (field.label === 'LABEL_REPEATED') {
      return (
        <div>
          {formFields.map((f, index) => (
            <div key={f.id}>
              <InputFormContents {...inputProps} index={index} />
              <Button onClick={() => remove(index)}>Remove</Button>
            </div>
          ))}
          <Button onClick={() => append({})}>Append</Button>
        </div>
      )
    }
    return <InputFormContents {...inputProps} />
  }

  const formattedFields: GrpcFormFieldType[] = [];
  desc.field.forEach((field) => {
    if (field.oneofIndex !== undefined) {
      const oneofType = desc.oneofDecl[field.oneofIndex]
      const existingOneof = formattedFields.find((f) => f.type === 'oneof' && f.name === oneofType.name);
      if (!existingOneof) {
        formattedFields.push({
          type: 'oneof',
          name: oneofType.name || 'unknown',
          fields: [field],
        })
      } else {
        if (existingOneof.type === 'oneof') {
          existingOneof.fields.push(field);
        }
      }
    } else {
      formattedFields.push({
        type: 'field',
        name: field.name || 'unknown',
        field,
      });
    }
  });
  return (
    <Accordion>
      {formattedFields.map((field) => {
        const panelContents = () => {
          if (field.type === 'field') {
            return formatField(field.field)
          } else {
            return (
              <Accordion>
                {field.fields.map((f) => (
                  <AccordionItem key={f.number} value={f.number}>
                    <AccordionHeader>{f.name}</AccordionHeader>
                    <AccordionPanel>{formatField(f)}</AccordionPanel>
                  </AccordionItem>
                ))}
              </Accordion>
            )
          }
        }
        return (
          <AccordionItem key={field.name} value={field.name}>
            <AccordionHeader>{field.name}</AccordionHeader>
            <AccordionPanel>{panelContents()}</AccordionPanel>
          </AccordionItem>
        )
      })}
    </Accordion>
  )
}

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
      input: node.data.config.grpc?.input,
      data: JSON.parse(localStorage.getItem(getNodeDataKey(node)) || '{}'),
    },
  });
  const values = watch();

  const onSubmit = (data: any) => {
    node.data.name = data.name;

    if (!node.data.config.grpc) {
      node.data.config.grpc = {
        package: "",
        service: "",
        method: "",
        // TODO breadchris this is not a valid grpc descriptor proto, but OK for now?
        // @ts-ignore
        input: {},
        data: {},
      };
    }

    node.data.config.grpc.package = data.config.package;
    node.data.config.grpc.service = data.config.service;
    node.data.config.grpc.method = data.config.method;

    localStorage.setItem(getNodeDataKey(node), JSON.stringify(data.data));
    onCancel();
  };

  const fieldPath = node.data.config.grpc.package;
  //@ts-ignore
  const inputFormProps: GRPCInputFormProps = {
    desc: node.data.config.grpc?.input,
    descLookup: node.data.config.grpc.descLookup,
    enumLookup: node.data.config.grpc.enumLookup,
    baseFieldName: 'data',
    //@ts-ignore
    register,
    // TODO breadchris without this ignore, my computer wants to take flight https://github.com/react-hook-form/react-hook-form/issues/6679
    //@ts-ignore
    control,
    fieldPath,
  }
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
        <GrpcInputForm {...inputFormProps} />
        <Divider/>
        <EditorActions/>
      </div>
    </form>
  );
}
