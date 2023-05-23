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
import {GetNodeInfoResponse, Project} from "@/rpc/project_pb";
import {useProjectResources} from "@/hooks/useProjectResources";
import {Resource} from "@/rpc/resource_pb";
import {toast} from "react-hot-toast";
import {useErrorBoundary} from "react-error-boundary";
import {getUpdatedProject} from "@/lib/project";
import {useEditorContext} from "@/providers/EditorProvider";

type ProjectContextType = {
    project: Project | undefined;
    resources: Resource[];
    output: any;
    loadingResources: boolean;

    saveProject: (nodes: Node[], edges: Edge[]) => Promise<void>;
    resetOutput: () => void;
    runWorkflow: (node: Node) => Promise<any>;
    loadResources: () => Promise<void>;
    deleteResource: (resourceId: string) => Promise<void>;
    loadNodeInfo: (nodeId: string) => Promise<GetNodeInfoResponse | undefined>;
};

type ProjectProviderProps = {
    children: React.ReactNode;
};

export function getNodeDataKey(node: Node) {
    return `${node.id}-sampleData`;
}

export function getDataFromNode(node: Node) {
    return localStorage.getItem(getNodeDataKey(node));
}

const ProjectContext = createContext<ProjectContextType>({} as any);

export const useProjectContext = () => useContext(ProjectContext);
export const useResetOutput = () => useProjectContext().resetOutput;

export default function ProjectProvider({children}: ProjectProviderProps) {
    const {project, loading, createDefault} = useDefaultProject();
    const {resources, loading: loadingResources, loadProjectResources} = useProjectResources();
    const [output, setOutput] = useState<any>(null);
    const {showBoundary} = useErrorBoundary();

    const resetOutput = useCallback(() => {
        setOutput(null);
    }, []);

    const saveProject = useCallback(async (nodes: Node[], edges: Edge[]) => {
        if (!project) return;

        const updatedProject = getUpdatedProject({
            project,
            nodes: nodes,
            edges: edges,
        });

        for (const node of updatedProject.graph?.nodes || []) {
            if (!node.name) {
                toast.error("Please name all nodes before exporting");
                return;
            }
        }

        await projectService.saveProject(updatedProject);
    }, [project]);

    const runWorkflow = useCallback(
        async (node: Node) => {
            if (!project) return;

            try {
                const res = await projectService.runWorklow({
                    nodeId: node.id,
                    projectId: project.id,
                    // TODO breadchris this is garbo, we need a better data structure to represent data input
                    input: getDataFromNode(node) || ''
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
        if (project) {
            void loadResources();
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
            }}
        >
            {children}
        </ProjectContext.Provider>
    );
}
