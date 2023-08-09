import React, {useCallback, useEffect, useState} from "react";
import { TextField, PrimaryButton, Stack, List, MessageBar } from '@fluentui/react';
import {Button, Card, SelectTabData, SelectTabEvent, Tab, TabList, TabValue} from "@fluentui/react-components";
import {useProjectContext} from "@/providers/ProjectProvider";
import {JsonViewer} from "@/components/jsonViewer";
import {toast} from "react-hot-toast";
import ReactMarkdown from "react-markdown";
import {useCurrentNode, useEditorContext} from "@/providers/EditorProvider";
import {GRPCInputFormProps, ProtobufInputForm} from "@/components/ProtobufInputForm";
import {useForm} from "react-hook-form";
import {Node as ProtoNode} from "@/rpc/graph_pb";

export function getNodeDataKey(node: ProtoNode) {
    return `${node.id}-sampleData`;
}

export function getDataFromNode(node: ProtoNode) {
    const data = localStorage.getItem(getNodeDataKey(node));
    return JSON.parse(data || '{}')
}

export function setDataForNode(node: ProtoNode, data: Object) {
    return localStorage.setItem(getNodeDataKey(node), JSON.stringify(data));
}

export default function RunPanel() {
    const {
        workflowOutput,
        setWorkflowOutput,
        project,
        runWorkflow,
    } = useProjectContext();
    const {nodeInfo} = useEditorContext();
    const currentNode = useCurrentNode();

    const {watch, setValue, register, handleSubmit, control} = useForm({
        values: {
            data: currentNode ? getDataFromNode(currentNode) : JSON.parse('{}'),
        },
    });

    const clearChat = () => {
        setWorkflowOutput([]);
    }

    const onRun = async (data: any) => {
        if (!project) return;
        if (!currentNode) {
            toast.error('Please select a node to run');
            return;
        }

        // TODO breadchris since there is no default value set for the form, it gets confused and by default
        // sets the fields to Array(0).
        const fix = (data: any) => {
            if (!data) {
                return data;
            }
            let newData: any = {};
            Object.keys(data).forEach((key) => {
                if (Array.isArray(data[key]) || typeof data[key] === 'object') {
                    newData[key] = fix(data[key])
                } else {
                    newData[key] = data[key];
                }
            })
            return newData;
        }

        data.data = fix(data.data);
        setDataForNode(currentNode, data.data);
        await runWorkflow(currentNode, data.data);
    };

    const onStartServer = async () => {
        if (!project) return;
        await runWorkflow(undefined, true);
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
        <Stack>
            {workflowOutput ? (
                <Stack>
                    <List items={workflowOutput} onRenderCell={(item?: string) => {
                        const getText = () => {
                            if (!item) {
                                return '';
                            }
                            const parsed = JSON.parse(item);
                            if (Object.keys(parsed).length === 1) {
                                const value: any = Object.values(parsed)[0];
                                return (<ReactMarkdown>{value.toString()}</ReactMarkdown>);
                            }
                            return <JsonViewer data={parsed} />;
                        }
                        return (
                            <MessageBar messageBarType={0}>{getText()}</MessageBar>
                        );
                    }} />
                    <Button onClick={clearChat}>Clear</Button>
                </Stack>
            ) : (
                <div>No results</div>
            ) }
            <form onSubmit={handleSubmit(onRun)}>
                {currentNode && form()}
                <Button onClick={onRun} type={'submit'}>
                    Run Workflow
                </Button>
            </form>
            <Button onClick={onStartServer}>
                Start Server
            </Button>
        </Stack>
    );
};
