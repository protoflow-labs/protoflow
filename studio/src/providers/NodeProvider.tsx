import {createContext, useContext, useEffect, useState} from "react";
import {GetNodeInfoResponse} from "@/rpc/project_pb";
import {useProjectContext} from "@/providers/ProjectProvider";

export type NodeContextType = {
    nodeId: string;
    nodeInfo: GetNodeInfoResponse | undefined;
}

const NodeContext = createContext<NodeContextType>({} as any);

export const useNodeContext = () => useContext(NodeContext);

export default function NodeProvider({children, nodeId}: { children: React.ReactNode, nodeId: string }) {
    const {loadNodeInfo} = useProjectContext();
    const [nodeInfo, setNodeInfo] = useState<GetNodeInfoResponse | undefined>(undefined);
    useEffect(() => {
        // TODO breadchris why is nodeId changing twice and triggering this twice?
        loadNodeInfo(nodeId).then(res => {
            setNodeInfo(res);
        });
    }, [nodeId]);

    return (
        <NodeContext.Provider value={{nodeId, nodeInfo}}>
            {children}
        </NodeContext.Provider>
    )
}