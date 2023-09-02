import React, {useState} from 'react';
import {makeStyles, shorthands, Switch} from "@fluentui/react-components";
import {List, PrimaryButton, Stack, TextField} from "@fluentui/react";
import {Toolbar} from "@/components/Chat/Toolbar";
import {Background, ReactFlow} from "reactflow";
import RunPanel from "@/components/Chat/RunPanel";
import {EditorPanel} from "@/components/EditorPanel/EditorPanel";
import {useActivelyEditing, useEditorContext} from "@/providers/EditorProvider";

interface Message {
    text: string;
    sender: 'user' | 'bot';
}

interface WindowProps {
    showingChat: boolean
}

export const Window: React.FC<WindowProps> = ({ showingChat }) => {
    const {props, setInstance} = useEditorContext();
    const [inputValue, setInputValue] = useState<string | undefined>('');

    const [messages, setMessages] = useState<Message[]>([
        {text: 'Hello, how can I help you?', sender: 'bot'},
    ]);

    const handleSend = () => {
        if (inputValue && inputValue.trim() !== '') {
            setMessages([...messages, {text: inputValue, sender: 'user'}]);
            setInputValue('');
        }
    };

    return (
        <>
            <Stack verticalFill verticalAlign="space-between"
                   styles={{root: {width: showingChat ? '80%' : '100%', margin: '0 auto', paddingTop: 10}}}>
                {showingChat ? (
                    <List
                        items={messages}
                        onRenderCell={(message) => {
                            if (!message) {
                                return null;
                            }
                            return (
                                <Stack horizontalAlign="stretch" styles={{root: {maxWidth: '100%'}}}>
                                    <div style={{
                                        backgroundColor: message.sender === 'user' ? '#e1e1e1' : '#0078d4',
                                        padding: 10,
                                        borderRadius: 5,
                                        color: message.sender === 'user' ? 'black' : 'white',
                                        alignSelf: message.sender === 'user' ? 'flex-end' : 'flex-start'
                                    }}>
                                        {message.text}
                                    </div>
                                </Stack>
                            )
                        }}
                    />
                ) : (
                    <div className="flex flex-1">
                        <Toolbar />
                        <ReactFlow
                            onInit={setInstance}
                            proOptions={{hideAttribution: true}}
                            {...props}
                            fitView
                        >
                            <Background/>
                        </ReactFlow>
                    </div>
                )}
                <Stack horizontal verticalAlign="end" horizontalAlign="center"
                       styles={{root: {width: '100%', gap: 15, marginBottom: 20, relative: true}}}>
                    {showingChat ? (
                        <>
                            <TextField
                                placeholder="Type a message..."
                                value={inputValue}
                                onChange={(e, newValue) => setInputValue(newValue)}
                                onKeyPress={(e) => {
                                    if (e.key === 'Enter') {
                                        handleSend();
                                    }
                                }}
                                underlined
                                styles={{root: {width: '100%'}}}
                            />
                            <PrimaryButton text="Send" onClick={handleSend}/>
                        </>
                    ) : (
                        <>
                            <RunPanel/>
                            <EditorPanel/>
                        </>
                    )}
                </Stack>
            </Stack>
        </>
    )
}
