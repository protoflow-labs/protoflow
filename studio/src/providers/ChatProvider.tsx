import {createContext, useContext, useEffect, useState} from "react";
import {GetNodeInfoResponse} from "@/rpc/project_pb";
import {useProjectContext} from "@/providers/ProjectProvider";

export type ChatContextType = {
}

const ChatContext = createContext<ChatContextType>({} as any);

export const useChatContext = () => useContext(ChatContext);

export default function ChatProvider({children}: { children: React.ReactNode }) {
    return (
        <ChatContext.Provider value={{}}>
            {children}
        </ChatContext.Provider>
    )
}
