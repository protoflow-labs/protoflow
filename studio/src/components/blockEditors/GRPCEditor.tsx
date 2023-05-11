import {
  Accordion,
  AccordionHeader,
  AccordionItem,
  AccordionPanel,
  Button,
  Divider,
  Field,
  Input,
} from "@fluentui/react-components";
import {Control, useFieldArray, useForm, UseFormRegister, useWatch} from "react-hook-form";
import {Node} from "reactflow";
import {GRPC} from "@/rpc/block_pb";
import {EditorActions, useUnselect} from "../EditorActions";
import {GRPCData} from "@/components/blocks/GRPCBlock";
import {DescriptorProto, FieldDescriptorProto, FieldDescriptorProto_Label} from "@bufbuild/protobuf";
import {FC, useState} from "react";

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
    return field.name;
  }
  if (idx !== undefined) {
    return `${baseFieldName}.${idx}.${field.name}`;
  }
  return `${baseFieldName}.${field.name}`;
}

type InputFormContentsProps = {
  baseFieldName: string,
  field: FieldDescriptorProto,
  descLookup: { [key: string]: DescriptorProto },
  register: UseFormRegister<any>,
  control: Control<GRPCFormValues, any>,
  index?: number,
}

const InputFormContents: FC<InputFormContentsProps> = (
  {
    baseFieldName,
    field,
    descLookup,
    register,
    control,
    index,
  }) => {
  const [visibleTypes, setVisibleTypes] = useState<string[]>([]);

  const fieldFormName = getFieldName(baseFieldName, field, index);
  console.log(fieldFormName)
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
                desc={fieldType}
                descLookup={descLookup}
                register={register}
                control={control}
                baseFieldName={fieldFormName}
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
  return (
    <Field key={field.number} label={field.name} required>
      <Input value={fieldValue} {...register(fieldFormName)} />
    </Field>
  )
}

const GrpcInputForm: FC<{
  desc: DescriptorProto | undefined,
  descLookup: { [key: string]: DescriptorProto },
  register: UseFormRegister<any>,
  control: Control<GRPCFormValues, any>,
  baseFieldName?: string,
}> = ({desc, descLookup, register, control, baseFieldName}) => {
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
      name: baseFieldName,
    });

    const props: InputFormContentsProps = {
      baseFieldName,
      field,
      register,
      descLookup,
      control,
    }

    // TODO breadchris for some reason FieldDescriptorProto_Label.REPEATED is a number and field.label is a string
    if (field.label === 'LABEL_REPEATED') {
      return (
        <div>
          {formFields.map((f, index) => (
            <div key={f.id}>
              <InputFormContents {...props} index={index} />
              <Button onClick={() => remove(index)}>Remove</Button>
            </div>
          ))}
          <Button onClick={() => append({})}>Append</Button>
        </div>
      )
    }
    return <InputFormContents {...props} />
  }
  const formattedFields: GrpcFormFieldType[] = [];
  desc.field.forEach((field) => {
    if (field.label === FieldDescriptorProto_Label.REPEATED) {

    }
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

type GRPCFormValues = {
  name: string
  config: {
    package: string
    service: string
    method: string
  },
  fields: GRPC
}

export function GRPCEditor({node}: { node: Node<GRPCData> }) {
  const onCancel = useUnselect();
  const {watch, setValue, register, handleSubmit, control} = useForm({
    values: {
      name: node.data.name || "",
      config: {
        package: node.data.config.grpc?.package || "",
        service: node.data.config.grpc?.service || "",
        method: node.data.config.grpc?.method || "",
      } as GRPC,
      fields: node.data.config.grpc?.input,
    },
  });

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
      };
    }

    node.data.config.grpc.package = data.config.package;
    node.data.config.grpc.service = data.config.service;
    node.data.config.grpc.method = data.config.method;

    onCancel();
  };

  const values = watch();
  console.log(values);

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
        {/* @ts-ignore */}
        <GrpcInputForm
          desc={node.data.config.grpc?.input}
          descLookup={node.data.config.grpc.descLookup}
          register={register}
          control={control}
          baseFieldName={'input'}
        />
        <Divider/>
        <EditorActions/>
      </div>
    </form>
  );
}
