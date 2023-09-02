import {DescriptorProto, FieldDescriptorProto, FieldDescriptorProto_Label} from "@bufbuild/protobuf";
import React, {FC} from "react";
import {useFieldArray} from "react-hook-form";
import {InputFormContents, InputFormContentsProps} from "@/components/ProtobufForm/InputFormContents";
import {Accordion, AccordionHeader, AccordionItem, AccordionPanel, Button} from "@fluentui/react-components";
import {GRPCInputFormProps} from "@/components/ProtobufForm/ProtobufInputForm";

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

export type ProtobufFormFieldType = GrpcFormField | GrpcFormOneof

interface GRPCInputFormContentsProps extends GRPCInputFormProps {
    field: ProtobufFormFieldType
    desc: DescriptorProto
}

export const AccordionField: FC<GRPCInputFormContentsProps> = (props) => {
    const {control, baseFieldName, field, fieldPath, desc} = props;

    const {fields: formFields, append, prepend, remove, swap, move, insert} = useFieldArray({
        control, name: baseFieldName || 'input',
    });

    const formatField = (field: FieldDescriptorProto) => {
        const inputProps: InputFormContentsProps = {
            ...props, field, fieldPath: `${fieldPath}.${desc.name}`
        }

        // TODO breadchris for some reason FieldDescriptorProto_Label.REPEATED is a number and field.label is a string
        // @ts-ignore
        if (field.label === 'LABEL_REPEATED' || field.label === FieldDescriptorProto_Label.REPEATED) {
            return (<div>
                    {formFields.map((f, index) => (<div key={f.id}>
                            <InputFormContents {...inputProps} index={index}/>
                            <Button onClick={() => remove(index)}>Remove</Button>
                        </div>))}
                    <Button onClick={() => append({})}>Append</Button>
                </div>)
        }
        return (<>
                <InputFormContents {...inputProps} />
            </>)
    }

    const panelContents = () => {
        if (field.type === 'field') {
            return formatField(field.field)
        } else {
            if (!field.fields) {
                return null;
            }
            return (<Accordion>
                    {field.fields.map((f) => (<AccordionItem key={f.name} value={f.number}>
                            <AccordionHeader>{f.name}</AccordionHeader>
                            <AccordionPanel>{formatField(f)}</AccordionPanel>
                        </AccordionItem>))}
                </Accordion>)
        }
    }
    return (<AccordionItem key={field.name} value={field.name}>
            <AccordionHeader>{field.name}</AccordionHeader>
            <AccordionPanel>{panelContents()}</AccordionPanel>
        </AccordionItem>)
}
