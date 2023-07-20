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

    const handleSendClick = async () => {
        if (inputValue) {
            const newMsgs = [...messages, new ChatMessage({
                role: 'user',
                message: inputValue
            })];
            setMessages(newMsgs);
            setInputValue('');

            const chat: Chat = new Chat({
                id: 'id',
            })
            const req: SendChatRequest = new SendChatRequest({
                chat: chat,
                message: inputValue,
            })
            let msg = '';
            try {
                for await (const message of projectService.sendChat(req)) {
                    msg += message.message;
                    setIncomingMessage((msg) => {
                        if (msg) {
                            return msg + message.message;
                        }
                        return message.message;
                    })
                }
            } catch (e) {
                console.error(e);
                return;
            }
            setIncomingMessage(null);
            setMessages([...newMsgs, new ChatMessage({
                role: 'bot',
                message: msg || 'No response'
            })]);
        }
    };

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
        <div className="absolute bottom-0 right-0 m-4 z-10 overflow-auto" style={{maxWidth: '400px', maxHeight: '500px'}}>
            <Card>
                <Stack horizontal={true}>
                    <TabList onTabSelect={tabSelect} vertical={true}>
                        <Tab value='chat'>Chat</Tab>
                        <Tab value='workflow'>Workflow</Tab>
                    </TabList>
                    <>
                        {activeTab === 'workflow' && (
                            <Stack>
                                {workflowOutput ? (
                                    <Stack>
                                        <List items={workflowOutput} onRenderCell={(item?: string) => {
                                            return (
                                                <MessageBar messageBarType={0} >
                                                    <JsonViewer data={item ? JSON.parse(item) : ''} />
                                                </MessageBar>
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
                        )}
                        {activeTab === 'chat' && (
                            <Stack>
                                <List items={messages} onRenderCell={(item?: ChatMessage) => <MessageBar messageBarType={0} >{item?.message}</MessageBar>} />
                                {incomingMessage && <MessageBar messageBarType={0} >{incomingMessage}</MessageBar>}
                                <Stack horizontal tokens={{ childrenGap: 10 }}>
                                    <TextField value={inputValue} onChange={(event, newValue) => {
                                        setInputValue(newValue || '');
                                    }} />
                                    <PrimaryButton text="Send" onClick={handleSendClick} />
                                </Stack>
                            </Stack>
                        )}
                    </>
                </Stack>
            </Card>
        </div>
    );
};
