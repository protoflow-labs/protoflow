import React from 'react'
import SyntaxHighlighter from "react-syntax-highlighter";
import {useEditorContext} from "@/providers/EditorProvider";

export const ProtoViewer: React.FC = ({}) => {
    const [open, setOpen] = React.useState(false);
    const { nodeInfo } = useEditorContext();
    return (
        <div>
            {(nodeInfo) && (
                <div className={"max-w-sm"}>
                    <SyntaxHighlighter language="protobuf">
                        {nodeInfo.methodProto}
                    </SyntaxHighlighter>
                </div>
            )}
        </div>
    )
}