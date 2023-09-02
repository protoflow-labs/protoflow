import React, {useState} from "react";
import {useProjectContext} from "@/providers/ProjectProvider";
import {useForm} from "react-hook-form";
import {GRPCInputFormProps, ProtobufInputForm} from "@/components/ProtobufForm/ProtobufInputForm";
import {toast} from "react-hot-toast";
import {Button, Divider, Textarea} from "@fluentui/react-components";
import {Edge as ProtoEdge} from "@/rpc/graph_pb";
import { GRPCTypeInfo } from "@/rpc/project_pb";

interface EdgeEditorProps {
    edge: ProtoEdge;
}

export const ActiveEdgeEditor: React.FC<EdgeEditorProps> = ({edge}) => {
    const { projectTypes , setEdgeLookup} = useProjectContext();
    const [config, setConfig] = useState<string>(JSON.stringify(edge.toJson()?.valueOf() || {}, null, 2));
    const { register, control, handleSubmit, setValue} = useForm({
        values: {
            data: edge.toJson()?.valueOf() || {}
        },
    });
    if (!projectTypes || !projectTypes.edgeType) {
        return null;
    }

    const inputFormProps: GRPCInputFormProps = {
        grpcInfo: new GRPCTypeInfo({
            input: projectTypes.edgeType,
            output: projectTypes.edgeType,
            descLookup: projectTypes.descLookup,
            enumLookup: projectTypes.enumLookup,
            packageName: '',
        }),
        // some random key to separate data from the form
        baseFieldName: 'data',
        //@ts-ignore
        register,
        setValue,
        // TODO breadchris without this ignore, my computer wants to take flight https://github.com/react-hook-form/react-hook-form/issues/6679
        //@ts-ignore
        control,
        fieldPath: '',
    }

    const onSubmit = async (data: any) => {
        setEdgeLookup((lookup) => {
            const edge = ProtoEdge.fromJson(data.data);
            return {
                ...lookup,
                [edge.id]: edge,
            }
        })
        toast.success('Saved!');
    };

    const saveConfig = () => {
        setEdgeLookup((lookup) => {
            const edge = ProtoEdge.fromJson(JSON.parse(config));
            return {
                ...lookup,
                [edge.id]: edge,
            }
        })
        toast.success('Saved!');
    }

    return (
        <div className="flex flex-col gap-2 p-3">
            <form onSubmit={handleSubmit(onSubmit)}>
                <div className="flex flex-col gap-2 p-3">
                    <ProtobufInputForm {...inputFormProps} />
                </div>
                <div className="flex items-center">
                    <Button appearance="primary" type="submit">
                        Save
                    </Button>
                </div>
                <Divider/>
                <Textarea value={config} onChange={(e) => setConfig(e.target.value)} />
                <Button onClick={saveConfig}>Save</Button>
            </form>
        </div>
    );
}

