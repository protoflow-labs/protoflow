import React, {useEffect, useState} from 'react';
import {Button, Switch, Tab, TabList} from "@fluentui/react-components";
import ProviderList from "@/components/Chat/ProviderList";
import {useProjectContext} from "@/providers/ProjectProvider";

interface SidebarProps {
    showingChat: boolean
    setShowingChat: (show: boolean) => void
}

export const Sidebar: React.FC<SidebarProps> = ({ showingChat, setShowingChat }) => {
    const [tab, setTab] = useState<string>('providers');
    const { runningWorkflows, loadRunningWorkflows, nodeLookup, stopWorkflow } = useProjectContext();
    const onShowingChatChange = React.useCallback(
        (ev: { currentTarget: { checked: boolean }; }) => {
            setShowingChat(ev.currentTarget.checked);
        },
        [setShowingChat]
    );

    useEffect(() => {
        if (tab === 'workflows') {
            void loadRunningWorkflows();
        }
    }, [loadRunningWorkflows, tab]);

    return (
        <div>
            <Switch
                checked={showingChat}
                onChange={onShowingChatChange}
                label={showingChat ? "Chat" : "Project"}
            />
            <TabList onTabSelect={(e, v) => {setTab(v.value as string)}}>
                <Tab value="providers">Providers</Tab>
                <Tab value="workflows">Workflows</Tab>
            </TabList>
            {tab === 'providers' ? <ProviderList/> : (
                <>
                    {runningWorkflows.map((workflow) => {
                        if (!workflow.request || !nodeLookup[workflow.request.nodeId]) {
                            return null;
                        }
                        const n = nodeLookup[workflow.request.nodeId];
                        return <div key={workflow.id}>Running {n.name}<Button onClick={() => stopWorkflow(workflow.id)}>Stop</Button></div>
                    })}
                </>
            )}
            {/*<TabList vertical size={"medium"}>*/}
            {/*    {previousChats.map((chat) => {*/}
            {/*        return <Tab key={chat.key} value={chat.key} icon={<Icon iconName={chat.icon} />}>{chat.name}</Tab>*/}
            {/*    })}*/}
            {/*</TabList>*/}
        </div>
    );
}