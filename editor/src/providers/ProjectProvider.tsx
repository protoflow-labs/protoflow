import { useDefaultProject } from "@/hooks/useDefaultProject";
import { projectService } from "@/lib/api";
import {
  Button,
  Caption1,
  Card,
  CardHeader,
  Spinner,
  Title3,
} from "@fluentui/react-components";
import { createContext, useCallback, useContext, useState } from "react";
import { HiExclamationCircle, HiPlus } from "react-icons/hi2";
import { Project } from "../../rpc/project_pb";

type ProjectContextType = {
  project: Project | undefined;
  output: any;

  resetOutput: () => void;
  runWorkflow: (nodeId: string) => Promise<any>;
  saveProject: () => Promise<void>;
};

type ProjectProviderProps = {
  children: React.ReactNode;
};

const ProjectContext = createContext<ProjectContextType>({} as any);

export const useProjectContext = () => useContext(ProjectContext);
export const useResetOutput = () => useProjectContext().resetOutput;

export default function ProjectProvider({ children }: ProjectProviderProps) {
  const { project, loading, createDefault } = useDefaultProject();
  const [output, setOutput] = useState<any>(null);

  const resetOutput = useCallback(() => {
    setOutput(null);
  }, []);

  const saveProject = useCallback(async () => {
    if (!project) return;

    await projectService.saveProject(project);
  }, [project]);

  const runWorkflow = useCallback(
    async (nodeId: string) => {
      if (!project) return;
      const res = await projectService.runWorklow({
        nodeId,
        projectId: project.id,
      });

      setOutput(res.output);
    },
    [project]
  );

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Spinner label="Loading project..." />
      </div>
    );
  }

  if (!project) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Card className="max-w-sm">
          <CardHeader
            image={<HiExclamationCircle className="text-blue-600 w-8 h-8" />}
            header={<Title3>No default project</Title3>}
            description={
              <Caption1>
                A default project could not be found. Use the button below to
                create one.
              </Caption1>
            }
          />
          <Button icon={<HiPlus />} onClick={createDefault}>
            Create Default
          </Button>
        </Card>
      </div>
    );
  }

  return (
    <ProjectContext.Provider
      value={{ project, output, resetOutput, runWorkflow, saveProject }}
    >
      {children}
    </ProjectContext.Provider>
  );
}
