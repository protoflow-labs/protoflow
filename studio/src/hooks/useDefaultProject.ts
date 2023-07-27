import { useEffect, useState } from "react";
import { Project, ProjectTypes } from "@/rpc/project_pb";
import { projectService } from "@/lib/api";

export function useDefaultProject() {
  const [loading, setLoading] = useState(true);
  const [project, setProject] = useState<Project>();
  const [projectTypes, setProjectTypes] = useState<ProjectTypes>();
  const loadProject = async () => {
    try {
      const { projects } = await projectService.getProjects({
        name: "local",
      });

      if (projects.length !== 1) {
        throw new Error(`No default project found: ${projects}`);
      }

      setProject(projects[0].project);
      setProjectTypes(projects[0].types);
    } catch (e) {
      console.error(e);
    }

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
    projectTypes,
    createDefault,
  };
}
