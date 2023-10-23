import {
    FieldDescriptorProto,
} from "@bufbuild/protobuf";
import React, {FC} from "react";
import {Control, UseFormRegister} from "react-hook-form";
import {
    Accordion,
} from "@fluentui/react-components";
import {GRPCTypeInfo} from "@/rpc/project_pb";
import {AccordionField, ProtobufFormFieldType} from "@/components/ProtobufFormSimple/AccordionField";

export interface GRPCInputFormProps {
    grpcInfo: GRPCTypeInfo
    register: UseFormRegister<any>
    control: Control
    fieldPath: string
    baseFieldName?: string
    setValue: (name: any, value: any) => void
}

export const ProtobufInputFormSimple: FC<GRPCInputFormProps> = (props) => {
    const {
        grpcInfo,
    } = props;

    const { input: desc } = grpcInfo;
    if (!desc) {
        return null;
    }

    const formattedFields: ProtobufFormFieldType[] = [];
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
        <ul className={"overflow-y-auto"} style={{maxHeight: "40em"}}>
            {formattedFields.map((field) => {
                return <AccordionField key={field.name} {...props} field={field} desc={desc} />
            })}
        </ul>
    )
}
