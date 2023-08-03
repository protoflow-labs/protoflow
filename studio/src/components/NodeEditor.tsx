import React from "react";
import {Node as ProtoNode} from "@/rpc/graph_pb";
import {useUnselect} from "@/components/EditorActions";
import {useEditorContext} from "@/providers/EditorProvider";
import {useNodeContext} from "@/providers/NodeProvider";
import {useForm} from "react-hook-form";
import {getNodeDataKey, useProjectContext} from "@/providers/ProjectProvider";
import {toast} from "react-hot-toast";
import {GRPCInputFormProps, ProtobufInputForm} from "@/components/ProtobufInputForm";
import {Button, Divider, Field, Input} from "@fluentui/react-components";

type NodeEditorProps = {
    node: ProtoNode;
};

export function NodeEditor(props: NodeEditorProps) {
    const {node} = props;

    const onCancel = useUnselect();
    const {save} = useEditorContext();
    const { setNodeLookup } = useProjectContext();
    const {nodeInfo} = useNodeContext();
    const {watch, setValue, register, handleSubmit, control} = useForm({
        values: {
            name: node.name || "",
            input: nodeInfo?.typeInfo?.input,
            data: node.toJson()?.valueOf(),
        },
    });
    const values = watch();

    const onSubmit = async (data: any) => {
        setNodeLookup((lookup) => {
            return {
                ...lookup,
                [node.id]: ProtoNode.fromJson(data.data),
            }
        })
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
