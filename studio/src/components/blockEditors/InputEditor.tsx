import {Button, Divider, Field, Input} from "@fluentui/react-components";
import React, {useEffect, useState} from "react";
import {EditorActions, useUnselect} from "../EditorActions";
import {ProtoViewer} from "@/components/blockEditors/ProtoViewer";
import {useNodeContext} from "@/providers/NodeProvider";
import {ProtobufInputForm, GRPCInputFormProps} from "@/components/inputForms/ProtobufInputForm";
import {useForm} from "react-hook-form";
import {getNodeDataKey} from "@/providers/ProjectProvider";
import {Node} from '@/rpc/graph_pb'
import {toast} from "react-hot-toast";
import {useEditorContext} from "@/providers/EditorProvider";


export function InputEditor({node}: { node: Node }) {
    const onCancel = useUnselect();
    const {save} = useEditorContext();
    const {nodeInfo} = useNodeContext();
    const {watch, setValue, register, handleSubmit, control} = useForm({
        values: {
            name: node.name || "",
            input: nodeInfo?.typeInfo?.input,
            data: JSON.parse(localStorage.getItem(getNodeDataKey(node)) || '{}'),
        },
    });
    const values = watch();

    const onSubmit = async (data: any) => {
        localStorage.setItem(getNodeDataKey(node), JSON.stringify(data.data));
        await save();
        toast.success('Saved!');
    };

    const form = () => {
        if (!nodeInfo || !nodeInfo.typeInfo) return (<></>);

        //@ts-ignore
        const inputFormProps: GRPCInputFormProps = {
            grpcInfo: nodeInfo.typeInfo,
            // some random key to separate data from the form
            baseFieldName: 'data',
            //@ts-ignore
            register,
            // TODO breadchris without this ignore, my computer wants to take flight https://github.com/react-hook-form/react-hook-form/issues/6679
            //@ts-ignore
            control,
            fieldPath: nodeInfo?.typeInfo?.packageName || '',
        }
        return (
            <ProtobufInputForm {...inputFormProps} />
        )
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="flex flex-col gap-2 p-3">
                <Field label="Name" required>
                    <Input value={values.name} {...register("name")} />
                </Field>
                <Divider/>
                {form()}
            </div>
            <div className="flex items-center">
                <Button appearance="primary" type="submit">
                    Save
                </Button>
            </div>
        </form>
    );
}
