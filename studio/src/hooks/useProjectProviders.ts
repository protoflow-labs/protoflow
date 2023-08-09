import { useEffect, useState } from "react";
import { projectService } from "@/lib/api";
import {EnumeratedProvider} from "@/rpc/project_pb";

export function useProjectProviders() {
  const [loading, setLoading] = useState(true);
  const [providers, setProviders] = useState<EnumeratedProvider[]>([]);
  const [providerLookup, setProviderLookup] = useState<Record<string, EnumeratedProvider>>({});
  const loadProjectResources = async (projectId: string) => {
    setLoading(true);
    try {
      const { providers } = await projectService.enumerateProviders({
        projectId,
      });
      setProviders(providers);
      setProviderLookup(providers.reduce((acc, provider) => {
        if (!provider.provider) {
          return acc;
        }
        return {
          ...acc,
          [provider.provider.id]: provider,
        };
      }, {} as Record<string, EnumeratedProvider>));
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  return {
    loading,
    providers: providers,
    providerLookup: providerLookup,
    loadProjectResources,
  };
}
