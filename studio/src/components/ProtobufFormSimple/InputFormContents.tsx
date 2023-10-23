import {FieldDescriptorProto, FieldDescriptorProto_Type} from "@bufbuild/protobuf";
import React, {FC, useState} from "react";
import {useWatch} from "react-hook-form";
import {Button, Field, Select, Textarea} from "@fluentui/react-components";
import {GRPCInputFormProps, ProtobufInputForm} from "@/components/ProtobufForm/ProtobufInputForm";
import { GRPCTypeInfo } from "@/rpc/project_pb";

const getFieldName = (baseFieldName: string | undefined, field: FieldDescriptorProto, idx?: number): string => {
    if (!baseFieldName) {
        return field.name || '';
    }
    if (idx !== undefined) {
        return `${baseFieldName}.${idx}.${field.name}`;
    }
    return `${baseFieldName}.${field.name}`;
}

export interface InputFormContentsProps extends GRPCInputFormProps {
    field: FieldDescriptorProto
    index?: number
}

export const InputFormContents: FC<InputFormContentsProps> = (props) => {
    const {
        grpcInfo,
        baseFieldName,
        field,
        register,
        control,
        index,
        fieldPath,
        setValue,
    } = props;
    const [visibleTypes, setVisibleTypes] = useState<string[]>([]);

    const { enumLookup, descLookup } = grpcInfo;

    const fieldFormName = getFieldName(baseFieldName, field, index);
    const fieldValue = useWatch({
        control,
        name: fieldFormName,
        defaultValue: '',
    });

    const useType = (typeName: string) => {
        //setValue(fieldFormName, {});
        setVisibleTypes([...visibleTypes, typeName]);
    }

    if (field.typeName) {
        // field.typeName == .name.othername, remove the leading dot
        const typeName = field.typeName.substring(1);
        const fieldType = descLookup[typeName];
        if (fieldType) {
            return (
                <div key={field.number}>
                    {visibleTypes.includes(typeName) ? (
                        <>
                            <ProtobufInputForm
                                {...props}
                                grpcInfo={new GRPCTypeInfo({
                                    ...grpcInfo,
                                    input: fieldType,
                                })}
                                baseFieldName={fieldFormName}
                                fieldPath={`${fieldPath}.${field.name}`}
                            />
                            <Button onClick={() => {
                                setVisibleTypes(visibleTypes.filter((t) => t !== typeName))
                            }}>Close</Button>
                        </>
                    ) : (
                        <Button onClick={() => useType(typeName)}>{typeName}</Button>
                    )}
                </div>
            )
        }
    }
    // TODO breadchris this should be checking for an actual type, not a string, field.type is a string, not a number
    // @ts-ignore
    if (field.type === "TYPE_ENUM" || field.type === FieldDescriptorProto_Type.ENUM) {
        if (!field.typeName) {
            throw new Error("Enum field has no type name");
        }
        const enumTypeName = `${fieldPath}.${field.name}`;
        const enumType = enumLookup[enumTypeName];
        console.log(enumLookup, enumTypeName, enumType)
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

