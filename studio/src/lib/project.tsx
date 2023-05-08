import { Edge, Node } from "reactflow";
import {Project, SaveProjectRequest} from "@/rpc/project_pb";
import {Node as ProtoNode, Edge as ProtoEdge, Graph} from "@/rpc/graph_pb";

export function getUpdatedProject({
  project,
  nodes,
  edges,
}: {
  project: Project;
  nodes: Node[];
  edges: Edge[];
}): SaveProjectRequest {
  return new SaveProjectRequest({
    projectId: project.id,
    graph: new Graph({
      id: project.graph?.id || project.id,
      name: project.graph?.name || project.name,
      edges: edges.map((edge) => (new ProtoEdge({
        id: edge.id,
        from: edge.source,
        to: edge.target,
      }))),
      nodes: nodes.map((node) => {
        const blockType: any = node.type?.split(".").pop();

        return new ProtoNode({
          id: node.id,
          name: node.data.name,

          x: node.position.x,
          y: node.position.y,
          config: {
            case: blockType,
            value: node.data.config[blockType] || node.data.config,
          },
          resourceIds: node.data.resourceIds || [],
        });
      }),
    }),
  });
}
