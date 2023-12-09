import { useProjectContext } from "@/providers/ProjectProvider";
import {PlugDisconnected20Regular} from "@fluentui/react-icons";
import React from "react";
import { ProviderState } from "@/rpc/project_pb";
import { NodeButton } from "@/components/Chat/ProviderList";
import {Manipulation} from "@/components/Builder/Manipulation";
import {Stack} from "@fluentui/react";
import {Textarea} from "@fluentui/react-components";
import {EditorPanel} from "@/components/EditorPanel/EditorPanel";
import {MethodInputForm} from "@/components/Builder/MethodInputForm";
import {useEditorContext} from "@/providers/EditorProvider";
import {ServiceSelector} from "@/components/Builder/ServiceSelector";

export function Builder() {
    const { project, providers } = useProjectContext();
    const { selectedNodes } = useEditorContext();

    if (!project) {
        return null;
    }

    console.log(selectedNodes)

    return (
        <main className="flex">
            <Stack>
                <Stack.Item>
                    <h1>Builder</h1>
                </Stack.Item>
                <Stack.Item>
                    <Stack horizontal>
                        <Stack.Item>
                            <ServiceSelector />
                        </Stack.Item>
                    </Stack>
                </Stack.Item>
            </Stack>
        </main>
    );
}
