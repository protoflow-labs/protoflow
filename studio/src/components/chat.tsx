import React, {useCallback, useEffect, useState} from "react";
import { TextField, PrimaryButton, Stack, List, MessageBar } from '@fluentui/react';
import {Button, Card, SelectTabData, SelectTabEvent, Tab, TabList, TabValue} from "@fluentui/react-components";
import {projectService} from "@/lib/api";
import { Chat, SendChatRequest, ChatMessage } from "@/rpc/project_pb";
import {useProjectContext} from "@/providers/ProjectProvider";
import {JsonViewer} from "@/components/jsonViewer";
import {useSelectedNodes} from "@/hooks/useSelectedNodes";
import {toast} from "react-hot-toast";
import {start} from "repl";
import ReactMarkdown from "react-markdown";

export default function ChatPanel() {
    const {
        workflowOutput,
        setWorkflowOutput,
        project,
        runWorkflow,
    } = useProjectContext();
    const { selectedNodes } = useSelectedNodes();
    const [inputValue, setInputValue] = useState('');
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [incomingMessage, setIncomingMessage] = useState<string | null>(null);
    const [activeTab, setActiveTab] = useState<TabValue>('workflow');


    const tabSelect = (event: SelectTabEvent, data: SelectTabData) => {
        setActiveTab(data.value);
    };

    const clearChat = () => {
        setWorkflowOutput([]);
    }

    const onRun = async () => {
        if (!project) return;
        if (!selectedNodes.length) {
            toast.error('Please select a node to run');
            return;
        }

        await runWorkflow(selectedNodes[0]);
    };

    const onStartServer = async () => {
        if (!project) return;
        await runWorkflow(undefined, true);
    };

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
            <Button onClick={onRun}>
                Run Workflow
            </Button>
            <Button onClick={onStartServer}>
                Start Server
            </Button>
        </Stack>
    );
};
