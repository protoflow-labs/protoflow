import React, {useCallback, useEffect, useState} from "react";
import { TextField, PrimaryButton, Stack, List, MessageBar } from '@fluentui/react';
import {Card} from "@fluentui/react-components";
import {projectService} from "@/lib/api";
import { Chat, SendChatRequest, ChatMessage } from "@/rpc/project_pb";

export default function ChatPanel() {
    const [inputValue, setInputValue] = useState('');
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [incomingMessage, setIncomingMessage] = useState<string | null>(null);

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

    return (
        <div className="absolute top-0 right-0 m-4 z-10 overflow-auto" style={{maxWidth: '400px'}}>
            <Card>
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
            </Card>
        </div>
    );
};
