import { useProjectContext } from "@/providers/ProjectProvider";
import {PlugDisconnected20Regular} from "@fluentui/react-icons";
import React from "react";
import { ProviderState } from "@/rpc/project_pb";
import { NodeButton } from "@/components/Chat/ProviderList";
import {Manipulation} from "@/components/Builder/Test";
import {Stack} from "@fluentui/react";
import {Textarea} from "@fluentui/react-components";
import {EditorPanel} from "@/components/EditorPanel/EditorPanel";
import {ActiveNodeEditor} from "@/components/EditorPanel/ActiveNodeEditor";
import {useEditorContext} from "@/providers/EditorProvider";

export function Builder() {
    const { project, providers } = useProjectContext();
    const { selectedNodes } = useEditorContext();

    if (!project) {
        return null;
    }

    return (
        <main className="flex">
            <Stack>
                <Stack.Item>
                    <h1>Builder</h1>
                </Stack.Item>
                <Stack.Item>
                    <Stack horizontal>
                        <Stack.Item>
                            <Manipulation />
                        </Stack.Item>
                        <Stack.Item>
                            {selectedNodes.map((node) => {
                                return (
                                    <ActiveNodeEditor key={node.id} node={node} />
                                )
                            })}
                        </Stack.Item>
                    </Stack>
                </Stack.Item>
            </Stack>
        </main>
    );
}
