import React from 'react'
import SyntaxHighlighter from "react-syntax-highlighter";
import {useNodeContext} from "@/providers/NodeProvider";
import {Button} from "@fluentui/react-components";

export const ProtoViewer: React.FC = ({}) => {
    const [open, setOpen] = React.useState(false);
    const {nodeInfo} = useNodeContext();
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