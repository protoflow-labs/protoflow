import {
    Badge,
    Divider,
    Field,
    Input,
    Textarea,
} from "@fluentui/react-components";
import {useForm} from "react-hook-form";
import {Node} from "@/rpc/graph_pb";
import {EditorActions, useUnselect} from "../EditorActions";
import {useProjectContext} from "@/providers/ProjectProvider";
import React from "react";
import {ProtoViewer} from "@/components/blockEditors/ProtoViewer";
import {FieldList, JsonReadOptions, JsonValue} from "@bufbuild/protobuf";

export interface NodeConfigType<T> {
    fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): T
    fields: FieldList
}

export function GenericNodeEditor<T>({node, nodeConfig, nodeConfigType}: { node: Node, nodeConfig: Node['config']['case'], nodeConfigType: NodeConfigType<T> }) {
    const onCancel = useUnselect();
    const {resources, setNodeLookup} = useProjectContext();
    if (node.config.case !== nodeConfig || !node.config.value) {
        return <div>Invalid node config</div>;
    }
    const {watch, setValue, register, handleSubmit, control} = useForm({
        values: {
            name: node.name || "",
            config: node.config.value.toJson(),
        },
    });
    const values = watch();

    const onSubmit = (data: any) => {
        if (!nodeConfig) {
            return;
        }
        node.name = data.name;
        node.config = {
            case: nodeConfig,
            // @ts-ignore
            value: nodeConfigType.fromJson(data.config),
        }
        setNodeLookup((prev) => {
            return {
                ...prev,
                [node.id]: node,
            }
        });
        onCancel();
    };

    const getResourceBadge = () => {
        if (!resources) {
            return null;
        }
        const res = resources.find((r) => {
            return r.resource && node.resourceId === r.resource.id
        })
        if (!res || !res.resource) {
            return null;
        }
        return <Badge key={res.resource.id}>{res.resource.name}</Badge>
    }

    const fields = nodeConfigType.fields.list();

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="flex flex-col gap-2 p-3">
                <Field label="Name" required>
                    <Input value={values.name} {...register("name")} />
                </Field>
                {fields.map((field) => {
                    return (
                        <Field label={field.name} key={field.name}>
                            {/* @ts-ignore */}
                            <Textarea value={values.config[field.name] || ''} {...register(`config.${field.name}`)} />
                        </Field>
                    )
                })}
                <Divider/>
                {getResourceBadge()}
                <Divider/>
                <ProtoViewer />
                <Divider/>
                <EditorActions/>
            </div>
        </form>
    );
}
