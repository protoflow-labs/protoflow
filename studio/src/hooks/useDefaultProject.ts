import { useEffect, useState } from "react";
import { Project } from "../../rpc/project_pb";
import { projectService } from "../lib/api";

export function useDefaultProject() {
  const [loading, setLoading] = useState(true);
  const [project, setProject] = useState<Project>();
  const loadProject = async () => {
    try {
      const { projects } = await projectService.getProjects({});

      const defaultProject = projects.find((p) => p.name === "local");
      if (!defaultProject) {
        throw new Error("No default project found");
      }

      setProject(defaultProject);
    } catch (e) {}

    setTimeout(() => {
      setLoading(false);
    }, 1);
  };

  const createDefault = async () => {
    const { project } = await projectService.createProject({
      name: "local",
    });

    setProject(project);
  };

  useEffect(() => {
    loadProject();
  }, []);

  return {
    loading,
    project,
    createDefault,
  };
}
