import {useDefaultProject} from "@/hooks/useDefaultProject";
import {projectService} from "@/lib/api";
import {
    Button,
    Caption1,
    Card,
    CardHeader,
    Spinner,
    Title3,
} from "@fluentui/react-components";
import {createContext, useCallback, useContext, useEffect, useState} from "react";
import {HiExclamationCircle, HiPlus} from "react-icons/hi2";
import {Edge, Node} from "reactflow";
import {EnumeratedProvider, GetNodeInfoResponse, Project, ProjectTypes} from "@/rpc/project_pb";
import {useProjectProviders} from "@/hooks/useProjectProviders";
import {toast} from "react-hot-toast";
import {useErrorBoundary} from "react-error-boundary";
import {getUpdatedProject} from "@/lib/project";
import {Node as ProtoNode, Edge as ProtoEdge} from "@/rpc/graph_pb";

type GetLookup = (lookup: Record<string, ProtoNode>) => Record<string, ProtoNode>;


type ProjectContextType = {
    project: Project | undefined;
    providers: EnumeratedProvider[];
    providerLookup: Record<string, EnumeratedProvider>;
    loadingResources: boolean;
    workflowOutput: string[] | null;
    setWorkflowOutput: (output: string[] | null) => void;

    saveProject: (nodes: Node[], edges: Edge[]) => Promise<void>;
    runWorkflow: (node?: Node, startServer?: boolean) => Promise<any>;
    loadResources: () => Promise<void>;
    loadNodeInfo: (nodeId: string) => Promise<GetNodeInfoResponse | undefined>;
    activeNode: ProtoNode | null;
    activeEdge: ProtoEdge | null;
    setActiveNodeId: (nodeId: string | null) => void;
    setActiveEdgeId: (edgeId: string | null) => void;

    projectTypes?: ProjectTypes;

    setNodeLookup: (getLookup: GetLookup) => void;
    nodeLookup: Record<string, ProtoNode>;
};

type ProjectProviderProps = {
    children: React.ReactNode;
};

export function getNodeDataKey(node: ProtoNode) {
    return `${node.id}-sampleData`;
}

export function getDataFromNode(node: ProtoNode) {
    return localStorage.getItem(getNodeDataKey(node));
}

const ProjectContext = createContext<ProjectContextType>({} as any);

export const useProjectContext = () => useContext(ProjectContext);

// project provider holds things that are closer to the database, like information fetched from the database
export default function ProjectProvider({children}: ProjectProviderProps) {
    const {project, loading, createDefault, projectTypes} = useDefaultProject();
    const {providers, providerLookup, loading: loadingResources, loadProjectResources} = useProjectProviders();
    const {showBoundary} = useErrorBoundary();
    const [nodeLookup, setNodeLookup] = useState<Record<string, ProtoNode>>({});
    const [edgeLookup, setEdgeLookup] = useState<Record<string, ProtoEdge>>({});
    const [activeNode, setActiveNode] = useState<ProtoNode | null>(null);
    const [activeEdge, setActiveEdge] = useState<ProtoEdge | null>(null);
    const [workflowOutput, setWorkflowOutput] = useState<string[] | null>(null);

    const setActiveNodeId = (nodeId: string | null) => {
        if (!nodeId) {
            setActiveNode(null);
            return;
        }

        // TODO breadchris catch error
        setActiveNode(nodeLookup[nodeId]);
    }

    const setActiveEdgeId = (edgeId: string | null) => {
        if (!edgeId) {
            setActiveEdge(null);
            return;
        }

        // TODO breadchris catch error
        setActiveEdge(edgeLookup[edgeId]);
    }

    const saveProject = useCallback(async (nodes: Node[], edges: Edge[]) => {
        if (!project) return;

        const updatedProject = getUpdatedProject(project, nodes, edges, nodeLookup);
        for (const node of updatedProject.graph?.nodes || []) {
            if (!node.name) {
                toast.error("Please name all nodes before exporting");
                return;
            }
        }

        await projectService.saveProject({
            projectId: project.id,
            graph: updatedProject.graph,
        });
    }, [nodeLookup]);

    const runWorkflow = useCallback(
        async (node?: Node, startServer?: boolean) => {
            if (!project) return;

            // TODO breadchris cleanup
            if (!node) {
                try {
                    setWorkflowOutput(null);
                    const res = await projectService.runWorkflow({
                        projectId: project.id,
                        startServer,
                    });
                    for await (const exec of res) {
                        setWorkflowOutput((prevState) => [...(prevState || []), exec.output]);
                    }
                } catch (e: any) {
                    toast.error(e.toString(), {
                        duration: 10000,
                    });
                }
                return;
            }

            const graphNode = nodeLookup[node.id];
            if (!graphNode) {
                toast.error(`Could not find node: ${node.id}`);
                return;
            }
            try {
                setWorkflowOutput(null);
                const res = await projectService.runWorkflow({
                    nodeId: node.id,
                    projectId: project.id,
                    // TODO breadchris this is garbo, we need a better data structure to represent data input
                    input: getDataFromNode(graphNode) || '',
                    startServer,
                });
                for await (const exec of res) {
                    setWorkflowOutput((prevState) => [...(prevState || []), exec.output]);
                }
            } catch (e: any) {
                toast.error(e.toString(), {
                    duration: 10000,
                });
            }
        },
        [project, nodeLookup]
    );

    const loadProviders = useCallback(async () => {
        if (!project) return;

        await loadProjectResources(project.id);
    }, [project]);

    const loadNodeInfo = useCallback(async (nodeId: string): Promise<GetNodeInfoResponse | undefined> => {
        if (!project) return undefined;

        try {
            return await projectService.getNodeInfo({
                nodeId: nodeId,
                projectId: project.id,
            });
        } catch (e) {
            // this is ok if we error, the node might not exist yet
            console.warn(e);
        }
        return undefined;
    }, [project]);

    // TODO breadchris should this happen every time the project is changed?
    useEffect(() => {
        void loadProviders();
        if (project && project.graph) {
            const lookup = project.graph.nodes.reduce((acc, node) => {
                acc[node.id] = node;
                return acc;
            }, {} as Record<string, ProtoNode>);

            const edgeLookup = project.graph.edges.reduce((acc, edge) => {
                acc[edge.id] = edge;
                return acc;
            }, {} as Record<string, ProtoEdge>);

            setNodeLookup((prev) => {
                return lookup;
            })

            setEdgeLookup((prev) => {
                return edgeLookup;
            });
        }
    }, [project]);

    if (loading) {
        return (
            <div className="flex items-center justify-center h-screen">
                <Spinner label="Loading project..."/>
            </div>
        );
    }

    if (!project) {
        return (
            <div className="flex items-center justify-center h-screen">
                <Card className="max-w-sm">
                    <CardHeader
                        image={<HiExclamationCircle className="text-blue-600 w-8 h-8"/>}
                        header={<Title3>No default project</Title3>}
                        description={
                            <Caption1>
                                A default project could not be found. Use the button below to
                                create one.
                            </Caption1>
                        }
                    />
                    <Button icon={<HiPlus/>} onClick={createDefault}>
                        Create Default
                    </Button>
                </Card>
            </div>
        );
    }

    return (
        <ProjectContext.Provider
            value={{
                project,
                providers,
                providerLookup,
                runWorkflow,
                workflowOutput,
                setWorkflowOutput,
                saveProject,
                loadResources: loadProviders,
                loadingResources,
                loadNodeInfo,
                setActiveNodeId,
                setActiveEdgeId,
                projectTypes,
                activeNode,
                activeEdge,
                setNodeLookup,
                nodeLookup,
            }}
        >
            {children}
        </ProjectContext.Provider>
    );
}
