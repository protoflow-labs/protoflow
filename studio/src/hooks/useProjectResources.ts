import { useEffect, useState } from "react";
import { projectService } from "@/lib/api";
import {Resource} from "@/rpc/resource_pb";

export function useProjectResources() {
  const [loading, setLoading] = useState(true);
  const [resources, setResources] = useState<Resource[]>([]);
  const loadProjectResources = async (projectId: string) => {
    setLoading(true);
    try {
      const { resources } = await projectService.getResources({
        projectId,
      });
      setResources(resources);
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  return {
    loading,
    resources,
    loadProjectResources,
  };
}