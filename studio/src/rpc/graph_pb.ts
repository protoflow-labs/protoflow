// @generated by protoc-gen-es v1.2.1 with parameter "target=ts"
// @generated from file graph.proto (package graph, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { Data } from "./data/data_pb.js";
import { Reason } from "./reason/reason_pb.js";
import { GRPC } from "./grpc/grpc_pb.js";
import { HTTP } from "./http/http_pb.js";
import { Storage } from "./storage/storage_pb.js";
import { Code } from "./code/code_pb.js";

/**
 * @generated from message graph.Graph
 */
export class Graph extends Message<Graph> {
  /**
   * @generated from field: repeated graph.Node nodes = 1;
   */
  nodes: Node[] = [];

  /**
   * @generated from field: repeated graph.Edge edges = 2;
   */
  edges: Edge[] = [];

  constructor(data?: PartialMessage<Graph>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "graph.Graph";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "nodes", kind: "message", T: Node, repeated: true },
    { no: 2, name: "edges", kind: "message", T: Edge, repeated: true },
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
   * @generated from oneof graph.Node.type
   */
  type: {
    /**
     * @generated from field: data.Data data = 7;
     */
    value: Data;
    case: "data";
  } | {
    /**
     * @generated from field: reason.Reason reason = 8;
     */
    value: Reason;
    case: "reason";
  } | {
    /**
     * @generated from field: grpc.GRPC grpc = 9;
     */
    value: GRPC;
    case: "grpc";
  } | {
    /**
     * @generated from field: http.HTTP http = 10;
     */
    value: HTTP;
    case: "http";
  } | {
    /**
     * @generated from field: storage.Storage storage = 11;
     */
    value: Storage;
    case: "storage";
  } | {
    /**
     * @generated from field: code.Code code = 12;
     */
    value: Code;
    case: "code";
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
    { no: 7, name: "data", kind: "message", T: Data, oneof: "type" },
    { no: 8, name: "reason", kind: "message", T: Reason, oneof: "type" },
    { no: 9, name: "grpc", kind: "message", T: GRPC, oneof: "type" },
    { no: 10, name: "http", kind: "message", T: HTTP, oneof: "type" },
    { no: 11, name: "storage", kind: "message", T: Storage, oneof: "type" },
    { no: 12, name: "code", kind: "message", T: Code, oneof: "type" },
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
 * @generated from message graph.Provides
 */
export class Provides extends Message<Provides> {
  constructor(data?: PartialMessage<Provides>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "graph.Provides";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Provides {
    return new Provides().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Provides {
    return new Provides().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Provides {
    return new Provides().fromJsonString(jsonString, options);
  }

  static equals(a: Provides | PlainMessage<Provides> | undefined, b: Provides | PlainMessage<Provides> | undefined): boolean {
    return proto3.util.equals(Provides, a, b);
  }
}

/**
 * @generated from message graph.PublishesTo
 */
export class PublishesTo extends Message<PublishesTo> {
  /**
   * @generated from field: string code_adapter = 1;
   */
  codeAdapter = "";

  constructor(data?: PartialMessage<PublishesTo>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "graph.PublishesTo";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "code_adapter", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PublishesTo {
    return new PublishesTo().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PublishesTo {
    return new PublishesTo().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PublishesTo {
    return new PublishesTo().fromJsonString(jsonString, options);
  }

  static equals(a: PublishesTo | PlainMessage<PublishesTo> | undefined, b: PublishesTo | PlainMessage<PublishesTo> | undefined): boolean {
    return proto3.util.equals(PublishesTo, a, b);
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

  /**
   * @generated from oneof graph.Edge.type
   */
  type: {
    /**
     * @generated from field: graph.Provides provides = 5;
     */
    value: Provides;
    case: "provides";
  } | {
    /**
     * @generated from field: graph.PublishesTo publishes_to = 6;
     */
    value: PublishesTo;
    case: "publishesTo";
  } | { case: undefined; value?: undefined } = { case: undefined };

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
    { no: 5, name: "provides", kind: "message", T: Provides, oneof: "type" },
    { no: 6, name: "publishes_to", kind: "message", T: PublishesTo, oneof: "type" },
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

