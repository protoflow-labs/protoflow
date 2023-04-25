import { Node } from "../../rpc/graph_pb";

export const configTypes =
  Node.fields.list().filter((f) => f.oneof?.name === "config") || [];
