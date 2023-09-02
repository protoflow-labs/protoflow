import React from "react";
import { Stack, List, MessageBar } from '@fluentui/react';
import {
    Button,
    Dialog, DialogActions,
    DialogBody,
    DialogContent,
    DialogSurface,
    DialogTitle,
    DialogTrigger
} from "@fluentui/react-components";
import {useProjectContext} from "@/providers/ProjectProvider";
import {JsonViewer} from "@/components/jsonViewer";
import {toast} from "react-hot-toast";
import ReactMarkdown from "react-markdown";
import {useCurrentNode, useEditorContext} from "@/providers/EditorProvider";
import {GRPCInputFormProps, ProtobufInputForm} from "@/components/ProtobufForm/ProtobufInputForm";
import {useForm} from "react-hook-form";
import {Node as ProtoNode} from "@/rpc/graph_pb";

export function getNodeDataKey(node: ProtoNode) {
    return `${node.id}-sampleData`;
}

export function getDataFromNode(node: ProtoNode) {
    const data = localStorage.getItem(getNodeDataKey(node));
    try {
        return JSON.parse(data || '{}')
    } catch (e) {
        console.log(e)
        return {};
    }
}

export function setDataForNode(node: ProtoNode, data: Object) {
    return localStorage.setItem(getNodeDataKey(node), JSON.stringify(data));
}

interface FormProps {
    data: any;
    onRun: (data: any) => void;
}

const Form: React.FC<FormProps> = ({ data, onRun }) => {
    const {nodeInfo} = useEditorContext();
    const {setValue, register, handleSubmit, control} = useForm({
        values: {
            data: data,
        },
    });

    if (!nodeInfo || !nodeInfo.typeInfo) {
        return null;
    }

    //@ts-ignore
    const inputFormProps: GRPCInputFormProps = {
        grpcInfo: nodeInfo.typeInfo,
        // some random key to separate data from the form
        baseFieldName: 'data',
        //@ts-ignore
        register,
        setValue,
        // TODO breadchris without this ignore, my computer wants to take flight https://github.com/react-hook-form/react-hook-form/issues/6679
        //@ts-ignore
        control,
        fieldPath: nodeInfo?.typeInfo?.packageName || '',
    }
    return (
        <>
            <DialogTitle>Dialog title</DialogTitle>
            <form onSubmit={handleSubmit(onRun)}>
                <DialogContent>
                    <ProtobufInputForm {...inputFormProps} />
                </DialogContent>
                <DialogActions>
                    <DialogTrigger disableButtonEnhancement>
                        <Button appearance="secondary">Close</Button>
                    </DialogTrigger>
                    <Button appearance="primary" type={'submit'}>Run</Button>
                </DialogActions>
            </form>
        </>
    );
}

export default function RunPanel() {
    const {
        workflowOutput,
        setWorkflowOutput,
        project,
        runWorkflow,
    } = useProjectContext();
    const currentNode = useCurrentNode();

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

    const justRun = async () => {
        if (!project) return;
        if (!currentNode) {
            toast.error('Please select a node to run');
            return;
        }
        await runWorkflow(currentNode, getDataFromNode(currentNode));
    }

    const clearInput = () => {
        if (!currentNode) {
            toast.error('Please select a node to run');
            return;
        }
        setDataForNode(currentNode, {});
    }

    const onStartServer = async () => {
        if (!project) return;
        await runWorkflow(undefined, {}, true);
    };

    const data = currentNode ? getDataFromNode(currentNode) : JSON.parse('{}');

    return (
        <Stack horizontal verticalAlign="end" horizontalAlign="center"
               styles={{root: {width: '100%', gap: 15, marginBottom: 20}}}>
            {workflowOutput ? (
                <Stack>
                    <List items={workflowOutput} onRenderCell={(item?: string) => {
                        const getText = () => {
                            if (!item) {
                                return '';
                            }
                            const parsed = JSON.parse(item);
                            const p = JSON.parse(parsed.output);
                            // if (Object.keys(p).length === 1) {
                            //     const value: any = Object.values(p)[0];
                            //     return (<ReactMarkdown>{value.toString()}</ReactMarkdown>);
                            // }
                            return <JsonViewer data={p} />;
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
            <Button onClick={clearInput}>
                Clear
            </Button>
            <Button onClick={justRun}>
                Run
            </Button>
            <Button onClick={onStartServer}>
                Start Server
            </Button>
            <Dialog>
                <DialogTrigger disableButtonEnhancement>
                    <Button>Input</Button>
                </DialogTrigger>
                <DialogSurface>
                    <DialogBody>
                        <Form data={data} onRun={onRun} />
                    </DialogBody>
                </DialogSurface>
            </Dialog>
        </Stack>
    );
};
