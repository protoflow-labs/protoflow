import React from "react";
import {Node as ProtoNode} from "@/rpc/graph_pb";
import {useEditorContext} from "@/providers/EditorProvider";
import {useForm} from "react-hook-form";
import {useProjectContext} from "@/providers/ProjectProvider";
import {toast} from "react-hot-toast";
import {Button} from "@fluentui/react-components";
import {GRPCInputFormProps, ProtoForm} from "@/components/ProtoForm/ProtoForm";
import { GRPCTypeInfo } from "@/rpc/project_pb";

type NodeEditorProps = {
    typeInfo: GRPCTypeInfo;
    onSubmit: (data: any) => void;
};

export function MethodInputForm(props: NodeEditorProps) {
    const {typeInfo, onSubmit} = props;

    const {resetField, register, handleSubmit, control, setValue} = useForm({
        values: {
            data: {},
        },
    });

    const form = () => {
        const inputFormProps: GRPCInputFormProps = {
            grpcInfo: typeInfo,
            // some random key to separate data from the form
            baseFieldName: 'data',
            //@ts-ignore
            register,
            setValue,
            resetField,
            // TODO breadchris without this ignore, my computer wants to take flight https://github.com/react-hook-form/react-hook-form/issues/6679
            //@ts-ignore
            control,
            fieldPath: typeInfo?.packageName || '',
        }
        return (
            <ProtoForm {...inputFormProps} />
        )
    }

    return (
        <form onSubmit={handleSubmit((data: any) => {
            onSubmit(data.data);
        })}>
            <div className="flex flex-col gap-2 p-3">
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
