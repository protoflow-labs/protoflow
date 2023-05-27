import { Edge, Node } from "reactflow";
import {Project, SaveProjectRequest} from "@/rpc/project_pb";
import {Node as ProtoNode, Edge as ProtoEdge, Graph} from "@/rpc/graph_pb";

export function getUpdatedProject(
  project: Project,
  nodes: Node[],
  edges: Edge[],
  nodeLookup: Record<string, ProtoNode>
): Project {
  return new Project({
    id: project.id,
    graph: new Graph({
      id: project.graph?.id || project.id,
      name: project.graph?.name || project.name,
      edges: edges.map((edge) => (new ProtoEdge({
        id: edge.id,
        from: edge.source,
        to: edge.target,
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
