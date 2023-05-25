import {FieldDescriptorProto} from "@bufbuild/protobuf";
import React, {FC, useState} from "react";
import {Control, useFieldArray, UseFormRegister, useWatch} from "react-hook-form";
import {
    Accordion,
    AccordionHeader,
    AccordionItem,
    AccordionPanel,
    Button,
    Field,
    Select,
    Textarea,
} from "@fluentui/react-components";
import {GRPCTypeInfo} from "@/rpc/project_pb";

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
        grpcInfo,
        baseFieldName,
        field,
        register,
        control,
        index,
        fieldPath,
    } = props;
    const [visibleTypes, setVisibleTypes] = useState<string[]>([]);

    const { enumLookup, descLookup } = grpcInfo;

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
                            <GRPCInputForm
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
            <Textarea value={fieldValue} {...register(fieldFormName)} resize={'vertical'} />
        </Field>
    )
}

export interface GRPCInputFormProps {
    grpcInfo: GRPCTypeInfo
    register: UseFormRegister<any>
    control: Control
    fieldPath: string
    baseFieldName?: string
}

export const GRPCInputForm: FC<GRPCInputFormProps> = (props) => {
    const {
        grpcInfo,
        control,
        baseFieldName,
        fieldPath
    } = props;

    const { input: desc } = grpcInfo;

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
    desc.field.forEach((field: FieldDescriptorProto) => {
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
