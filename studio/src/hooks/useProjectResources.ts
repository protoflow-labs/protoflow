import { useEffect, useState } from "react";
import { projectService } from "@/lib/api";
import {Resource} from "@/rpc/resource_pb";
import {EnumeratedResource} from "@/rpc/project_pb";

export function useProjectResources() {
  const [loading, setLoading] = useState(true);
  const [resources, setResources] = useState<EnumeratedResource[]>([]);
  const [resourceLookup, setResourceLookup] = useState<Record<string, EnumeratedResource>>({});
  const loadProjectResources = async (projectId: string) => {
    setLoading(true);
    try {
      const { resources } = await projectService.getResources({
        projectId,
      });
      setResources(resources);
      setResourceLookup(resources.reduce((acc, resource) => {
        if (!resource.resource) {
          return acc;
        }
        return {
          ...acc,
          [resource.resource.id]: resource,
        };
      }, {} as Record<string, EnumeratedResource>));
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  return {
    loading,
    resources,
    resourceLookup,
    loadProjectResources,
  };
}
