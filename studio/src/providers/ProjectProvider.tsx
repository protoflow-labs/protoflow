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
import {EnumeratedResource, GetNodeInfoResponse, Project} from "@/rpc/project_pb";
import {useProjectResources} from "@/hooks/useProjectResources";
import {toast} from "react-hot-toast";
import {useErrorBoundary} from "react-error-boundary";
import {getUpdatedProject} from "@/lib/project";
import {Node as ProtoNode} from "@/rpc/graph_pb";
import {notEmpty} from "@/util/predicates";

type GetLookup = (lookup: Record<string, ProtoNode>) => Record<string, ProtoNode>;


type ProjectContextType = {
    project: Project | undefined;
    resources: EnumeratedResource[];
    output: any;
    loadingResources: boolean;

    saveProject: (nodes: Node[], edges: Edge[]) => Promise<void>;
    resetOutput: () => void;
    runWorkflow: (node: Node) => Promise<any>;
    loadResources: () => Promise<void>;
    deleteResource: (resourceId: string) => Promise<void>;
    loadNodeInfo: (nodeId: string) => Promise<GetNodeInfoResponse | undefined>;
    activeNode: ProtoNode | null;
    setActiveNodeId: (nodeId: string | null) => void;

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
export const useResetOutput = () => useProjectContext().resetOutput;

// project provider holds things that are closer to the database, like information fetched from the database
export default function ProjectProvider({children}: ProjectProviderProps) {
    const {project, loading, createDefault} = useDefaultProject();
    const {resources, loading: loadingResources, loadProjectResources} = useProjectResources();
    const [output, setOutput] = useState<any>(null);
    const {showBoundary} = useErrorBoundary();
    const [nodeLookup, setNodeLookup] = useState<Record<string, ProtoNode>>({});
    const [activeNode, setActiveNode] = useState<ProtoNode | null>(null);

    const setActiveNodeId = (nodeId: string | null) => {
        if (!nodeId) return;

        // TODO breadchris catch error
        setActiveNode(nodeLookup[nodeId]);
    }

    const resetOutput = useCallback(() => {
        setOutput(null);
    }, []);

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
            resources: resources.map(r => r.resource).filter(notEmpty),
        });
    }, [nodeLookup, resources]);

    const runWorkflow = useCallback(
        async (node: Node) => {
            if (!project) return;

            const graphNode = nodeLookup[node.id];
            if (!graphNode) {
                toast.error(`Could not find node: ${node.id}`);
                return;
            }
            try {
                const res = await projectService.runWorklow({
                    nodeId: node.id,
                    projectId: project.id,
                    // TODO breadchris this is garbo, we need a better data structure to represent data input
                    input: getDataFromNode(graphNode) || ''
                });

                setOutput(res.output);
            } catch (e) {
                // @ts-ignore
                toast.error(e.toString());
            }
        },
        [project]
    );

    const deleteResource = useCallback(
        async (resourceId: string) => {
            if (!project) return;

            try {
                const res = await projectService.deleteResource({
                    projectId: project.id,
                    resourceId,
                });
                toast.success('deleted resource');
            } catch (e) {
                // @ts-ignore
                toast.error(e.toString());
            }
        },
        [project]
    );

    const loadResources = useCallback(async () => {
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
        }
        return undefined;
    }, [project]);

    // TODO breadchris should this happen every time the project is changed?
    useEffect(() => {
        void loadResources();
        if (project && project.graph) {
            const lookup = project.graph.nodes.reduce((acc, node) => {
                acc[node.id] = node;
                return acc;
            }, {} as Record<string, ProtoNode>);

            setNodeLookup((prev) => {
                return lookup;
            })
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
                resources,
                output,
                resetOutput,
                runWorkflow,
                saveProject,
                deleteResource,
                loadResources,
                loadingResources,
                loadNodeInfo,
                setActiveNodeId,
                activeNode,
                setNodeLookup,
                nodeLookup,
            }}
        >
            {children}
        </ProjectContext.Provider>
    );
}
