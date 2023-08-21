import {useActivelyEditing, useEditorMode} from "@/providers/EditorProvider";
import {EditorPanel} from "./EditorPanel";
import React, {useEffect} from "react";
import {Card, SelectTabData, SelectTabEvent, Tab, TabList, TabValue} from "@fluentui/react-components";
import {Stack} from "@fluentui/react";
import RunPanel from "@/components/RunPanel";

export function ActionPanel() {
    const activelyEditing = useActivelyEditing();
    const [activeTab, setActiveTab] = React.useState<TabValue>('chat');
    const tabSelect = (event: SelectTabEvent, data: SelectTabData) => {
        setActiveTab(data.value);
    };

    useEffect(() => {
        setActiveTab('run')
    }, [activelyEditing]);

    return (
        <div className="absolute top-0 right-0 m-4 z-10 overflow-auto h-screen">
            <Card className="h-screen max-w-sm">
                <Stack>
                    <TabList onTabSelect={tabSelect} vertical={true}>
                        <Tab value='run'>Run</Tab>
                        <Tab value='edit'>Edit</Tab>
                    </TabList>
                    {activeTab === 'run' ? <RunPanel/> : <EditorPanel/>}
                </Stack>
            </Card>
        </div>
    );
}
