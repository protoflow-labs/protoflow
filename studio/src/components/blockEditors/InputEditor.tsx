import {InputData} from "@/components/blocks/InputBlock";
import {Divider, Field, Input} from "@fluentui/react-components";
import React, {useEffect, useState} from "react";
import {Node} from "reactflow";
import {EditorActions, useUnselect} from "../EditorActions";
import {ProtoViewer} from "@/components/blockEditors/common/ProtoViewer";
import {useNodeContext} from "@/providers/NodeProvider";
import {GRPCInputForm, GRPCInputFormProps} from "@/components/inputForms/GRPCInputForm";
import {useForm} from "react-hook-form";
import {getNodeDataKey} from "@/providers/ProjectProvider";


export function InputEditor({node}: { node: Node<InputData> }) {
    const onCancel = useUnselect();
    const {nodeInfo} = useNodeContext();
    const {watch, setValue, register, handleSubmit, control} = useForm({
        values: {
            name: node.data.name || "",
            input: nodeInfo?.typeInfo?.input,
            data: JSON.parse(localStorage.getItem(getNodeDataKey(node)) || '{}'),
        },
    });
    const values = watch();

    const onSubmit = (data: any) => {
        localStorage.setItem(getNodeDataKey(node), JSON.stringify(data.data));
        onCancel();
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
            <GRPCInputForm {...inputFormProps} />
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
                <Divider/>
                <ProtoViewer/>
                <Divider/>
                <EditorActions/>
            </div>
        </form>
    );
}
