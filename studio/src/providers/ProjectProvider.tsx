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
import {Edge, Node, useOnSelectionChange} from "reactflow";
import {EnumeratedProvider, GetNodeInfoResponse, Project, ProjectTypes, WorkflowTrace} from "@/rpc/project_pb";
import {useProjectProviders} from "@/hooks/useProjectProviders";
import {toast} from "react-hot-toast";
import {useErrorBoundary} from "react-error-boundary";
import {Node as ProtoNode, Edge as ProtoEdge, Graph} from "@/rpc/graph_pb";

type GetLookup = (lookup: Record<string, ProtoNode>) => Record<string, ProtoNode>;
type GetEdgeLookup = (lookup: Record<string, ProtoEdge>) => Record<string, ProtoEdge>;

type ProjectContextType = {
    project: Project | undefined;
    providers: EnumeratedProvider[];
    providerLookup: Record<string, EnumeratedProvider>;
    loadingResources: boolean;
    workflowOutput: string[] | null;
    setWorkflowOutput: (output: string[] | null) => void;

    saveProject: (nodes: Node[], edges: Edge[]) => Promise<void>;
    runWorkflow: (node?: ProtoNode, input?: Object, startServer?: boolean) => Promise<any>;
    loadProviders: () => Promise<void>;

    projectTypes?: ProjectTypes;

    setNodeLookup: (getLookup: GetLookup) => void;
    nodeLookup: Record<string, ProtoNode>;

    setEdgeLookup: (getLookup: GetEdgeLookup) => void;
    edgeLookup: Record<string, ProtoEdge>;

    runningWorkflows: WorkflowTrace[];
    loadRunningWorkflows: () => Promise<void>;
    stopWorkflow: (workflowID: string) => Promise<void>;
};

type ProjectProviderProps = {
    children: React.ReactNode;
};

const ProjectContext = createContext<ProjectContextType>({} as any);

export const useProjectContext = () => useContext(ProjectContext);

function getUpdatedProject(
    project: Project,
    nodes: Node[],
    edges: Edge[],
    nodeLookup: Record<string, ProtoNode>,
    edgeLookup: Record<string, ProtoEdge>
): Project {
    return new Project({
        id: project.id,
        graph: new Graph({
            edges: edges.filter(e => edgeLookup[e.id] !== undefined).map((edge) => (new ProtoEdge({
                ...edgeLookup[edge.id],
            }))),
            nodes: nodes.map((node) => {
                return {
                    ...nodeLookup[node.id],
                    x: node.position.x,
                    y: node.position.y
                }
            }),
        }),
    });
}

// project provider holds things that are closer to the database, like information fetched from the database
export default function ProjectProvider({children}: ProjectProviderProps) {
    const {project, loading, createDefault, projectTypes} = useDefaultProject();
    const {providers, providerLookup, loading: loadingResources, loadProjectResources} = useProjectProviders();
    const {showBoundary} = useErrorBoundary();
    const [nodeLookup, setNodeLookup] = useState<Record<string, ProtoNode>>({});
    const [edgeLookup, setEdgeLookup] = useState<Record<string, ProtoEdge>>({});
    const [workflowOutput, setWorkflowOutput] = useState<string[] | null>(null);
    const [runningWorkflows, setRunningWorkflows] = useState<WorkflowTrace[]>([]);

    const loadRunningWorkflows = useCallback(async () => {
        try {
            const { traces } = await projectService.getRunningWorkflows({});
            setRunningWorkflows(traces);
        } catch (e) {
            console.error(e);
        }
    }, [setRunningWorkflows]);

    const stopWorkflow = useCallback(async (workflowID: string) => {
        try {
            await projectService.stopWorkflow({ workflowId: workflowID });
            await loadRunningWorkflows();
        } catch (e) {
            console.error(e);
        }
    }, [loadRunningWorkflows]);

    const saveProject = useCallback(async (nodes: Node[], edges: Edge[]) => {
        if (!project) return;

        const updatedProject = getUpdatedProject(project, nodes, edges, nodeLookup, edgeLookup);
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
    }, [nodeLookup, edgeLookup]);

    const runWorkflow = useCallback(
        async (node?: ProtoNode, input?: Object, startServer?: boolean) => {
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

            console.log('executing node', node.name, input)

            try {
                setWorkflowOutput(null);
                const res = await projectService.runWorkflow({
                    nodeId: node.id,
                    projectId: project.id,
                    // TODO breadchris this is garbo, we need a better data structure to represent data input
                    input: JSON.stringify(input) || '',
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
                loadProviders,
                loadingResources,
                projectTypes,
                setNodeLookup,
                nodeLookup,
                setEdgeLookup,
                edgeLookup,
                runningWorkflows,
                loadRunningWorkflows,
                stopWorkflow,
            }}
        >
            {children}
        </ProjectContext.Provider>
    );
}
