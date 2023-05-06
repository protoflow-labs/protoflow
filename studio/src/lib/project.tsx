import { Edge, Node } from "reactflow";
import { Project } from "@/rpc/project_pb";

export function getUpdatedProject({
  project,
  nodes,
  edges,
}: {
  project: Project;
  nodes: Node[];
  edges: Edge[];
}) {
  return {
    projectId: project.id,
    graph: {
      id: project.graph?.id || project.id,
      name: project.graph?.name || project.name,
      edges: edges.map((edge) => ({
        id: edge.id,
        from: edge.source,
        to: edge.target,
      })),
      nodes: nodes.map((node) => {
        const blockType: any = node.type?.split(".").pop();

        return {
          id: node.id,
          name: node.data.name,

          x: node.position.x,
          y: node.position.y,
          config: {
            case: blockType,
            value: node.data.config[blockType] || node.data.config,
          },
        };
      }),
    },
  };
}
