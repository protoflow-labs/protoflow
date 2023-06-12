// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file graph.proto (package graph, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { Bucket, Collection, Function, GRPC, Input, Query, REST } from "./block_pb.js";

/**
 * @generated from message graph.Graph
 */
export class Graph extends Message<Graph> {
  /**
   * TODO breadchris get rid of id and name, they are not needed
   *
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: repeated graph.Node nodes = 3;
   */
  nodes: Node[] = [];

  /**
   * @generated from field: repeated graph.Edge edges = 4;
   */
  edges: Edge[] = [];

  constructor(data?: PartialMessage<Graph>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "graph.Graph";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "nodes", kind: "message", T: Node, repeated: true },
    { no: 4, name: "edges", kind: "message", T: Edge, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Graph {
    return new Graph().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Graph {
    return new Graph().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Graph {
    return new Graph().fromJsonString(jsonString, options);
  }

  static equals(a: Graph | PlainMessage<Graph> | undefined, b: Graph | PlainMessage<Graph> | undefined): boolean {
    return proto3.util.equals(Graph, a, b);
  }
}

/**
 * @generated from message graph.Node
 */
export class Node extends Message<Node> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: float x = 4;
   */
  x = 0;

  /**
   * @generated from field: float y = 5;
   */
  y = 0;

  /**
   * @generated from field: string resource_id = 6;
   */
  resourceId = "";

  /**
   * @generated from oneof graph.Node.config
   */
  config: {
    /**
     * @generated from field: block.REST rest = 7;
     */
    value: REST;
    case: "rest";
  } | {
    /**
     * @generated from field: block.GRPC grpc = 8;
     */
    value: GRPC;
    case: "grpc";
  } | {
    /**
     * @generated from field: block.Collection collection = 9;
     */
    value: Collection;
    case: "collection";
  } | {
    /**
     * @generated from field: block.Bucket bucket = 10;
     */
    value: Bucket;
    case: "bucket";
  } | {
    /**
     * @generated from field: block.Input input = 11;
     */
    value: Input;
    case: "input";
  } | {
    /**
     * @generated from field: block.Function function = 12;
     */
    value: Function;
    case: "function";
  } | {
    /**
     * @generated from field: block.Query query = 13;
     */
    value: Query;
    case: "query";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Node>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "graph.Node";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "x", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
    { no: 5, name: "y", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
    { no: 6, name: "resource_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 7, name: "rest", kind: "message", T: REST, oneof: "config" },
    { no: 8, name: "grpc", kind: "message", T: GRPC, oneof: "config" },
    { no: 9, name: "collection", kind: "message", T: Collection, oneof: "config" },
    { no: 10, name: "bucket", kind: "message", T: Bucket, oneof: "config" },
    { no: 11, name: "input", kind: "message", T: Input, oneof: "config" },
    { no: 12, name: "function", kind: "message", T: Function, oneof: "config" },
    { no: 13, name: "query", kind: "message", T: Query, oneof: "config" },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Node {
    return new Node().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Node {
    return new Node().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Node {
    return new Node().fromJsonString(jsonString, options);
  }

  static equals(a: Node | PlainMessage<Node> | undefined, b: Node | PlainMessage<Node> | undefined): boolean {
    return proto3.util.equals(Node, a, b);
  }
}

/**
 * @generated from message graph.Edge
 */
export class Edge extends Message<Edge> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string from = 2;
   */
  from = "";

  /**
   * @generated from field: string to = 3;
   */
  to = "";

  constructor(data?: PartialMessage<Edge>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "graph.Edge";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "from", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "to", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Edge {
    return new Edge().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Edge {
    return new Edge().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Edge {
    return new Edge().fromJsonString(jsonString, options);
  }

  static equals(a: Edge | PlainMessage<Edge> | undefined, b: Edge | PlainMessage<Edge> | undefined): boolean {
    return proto3.util.equals(Edge, a, b);
  }
}

