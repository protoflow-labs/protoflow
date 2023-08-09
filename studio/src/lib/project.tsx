import { Edge, Node } from "reactflow";
import {Project, SaveProjectRequest} from "@/rpc/project_pb";
import {Node as ProtoNode, Edge as ProtoEdge, Graph} from "@/rpc/graph_pb";

export function getUpdatedProject(
  project: Project,
  nodes: Node[],
  edges: Edge[],
  nodeLookup: Record<string, ProtoNode>,
  edgeLookup: Record<string, ProtoEdge>
): Project {
  return new Project({
    id: project.id,
    graph: new Graph({
      edges: edges.filter(e => edgeLookup[e.id] !== undefined).map((edge) => (new ProtoEdge({
        ...edgeLookup[edge.id],
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
