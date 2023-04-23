// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file workflow.proto (package workflow, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message workflow.ID
 */
export class ID extends Message<ID> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  constructor(data?: PartialMessage<ID>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.ID";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ID {
    return new ID().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ID {
    return new ID().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ID {
    return new ID().fromJsonString(jsonString, options);
  }

  static equals(a: ID | PlainMessage<ID> | undefined, b: ID | PlainMessage<ID> | undefined): boolean {
    return proto3.util.equals(ID, a, b);
  }
}

/**
 * @generated from message workflow.Run
 */
export class Run extends Message<Run> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  constructor(data?: PartialMessage<Run>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.Run";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Run {
    return new Run().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Run {
    return new Run().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Run {
    return new Run().fromJsonString(jsonString, options);
  }

  static equals(a: Run | PlainMessage<Run> | undefined, b: Run | PlainMessage<Run> | undefined): boolean {
    return proto3.util.equals(Run, a, b);
  }
}

/**
 * @generated from message workflow.WorkflowEntrypoint
 */
export class WorkflowEntrypoint extends Message<WorkflowEntrypoint> {
  /**
   * @generated from field: string workflowId = 1;
   */
  workflowId = "";

  /**
   * @generated from field: string nodeId = 2;
   */
  nodeId = "";

  constructor(data?: PartialMessage<WorkflowEntrypoint>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.WorkflowEntrypoint";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "workflowId", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "nodeId", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): WorkflowEntrypoint {
    return new WorkflowEntrypoint().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): WorkflowEntrypoint {
    return new WorkflowEntrypoint().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): WorkflowEntrypoint {
    return new WorkflowEntrypoint().fromJsonString(jsonString, options);
  }

  static equals(a: WorkflowEntrypoint | PlainMessage<WorkflowEntrypoint> | undefined, b: WorkflowEntrypoint | PlainMessage<WorkflowEntrypoint> | undefined): boolean {
    return proto3.util.equals(WorkflowEntrypoint, a, b);
  }
}

/**
 * @generated from message workflow.Workflow
 */
export class Workflow extends Message<Workflow> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: repeated workflow.Node nodes = 3;
   */
  nodes: Node[] = [];

  /**
   * @generated from field: repeated workflow.Edge edges = 4;
   */
  edges: Edge[] = [];

  /**
   * @generated from field: repeated workflow.Resource resources = 5;
   */
  resources: Resource[] = [];

  constructor(data?: PartialMessage<Workflow>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.Workflow";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "nodes", kind: "message", T: Node, repeated: true },
    { no: 4, name: "edges", kind: "message", T: Edge, repeated: true },
    { no: 5, name: "resources", kind: "message", T: Resource, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Workflow {
    return new Workflow().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Workflow {
    return new Workflow().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Workflow {
    return new Workflow().fromJsonString(jsonString, options);
  }

  static equals(a: Workflow | PlainMessage<Workflow> | undefined, b: Workflow | PlainMessage<Workflow> | undefined): boolean {
    return proto3.util.equals(Workflow, a, b);
  }
}

/**
 * @generated from message workflow.Resource
 */
export class Resource extends Message<Resource> {
  /**
   * @generated from oneof workflow.Resource.type
   */
  type: {
    /**
     * @generated from field: workflow.DBResource db = 1;
     */
    value: DBResource;
    case: "db";
  } | {
    /**
     * @generated from field: workflow.DocStoreResource docstore = 2;
     */
    value: DocStoreResource;
    case: "docstore";
  } | {
    /**
     * @generated from field: workflow.BucketResource bucket = 3;
     */
    value: BucketResource;
    case: "bucket";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Resource>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.Resource";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "db", kind: "message", T: DBResource, oneof: "type" },
    { no: 2, name: "docstore", kind: "message", T: DocStoreResource, oneof: "type" },
    { no: 3, name: "bucket", kind: "message", T: BucketResource, oneof: "type" },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Resource {
    return new Resource().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Resource {
    return new Resource().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Resource {
    return new Resource().fromJsonString(jsonString, options);
  }

  static equals(a: Resource | PlainMessage<Resource> | undefined, b: Resource | PlainMessage<Resource> | undefined): boolean {
    return proto3.util.equals(Resource, a, b);
  }
}

/**
 * @generated from message workflow.Function
 */
export class Function extends Message<Function> {
  /**
   * @generated from oneof workflow.Function.type
   */
  type: {
    /**
     * @generated from field: workflow.CodeFunction code = 1;
     */
    value: CodeFunction;
    case: "code";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Function>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.Function";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "code", kind: "message", T: CodeFunction, oneof: "type" },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Function {
    return new Function().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Function {
    return new Function().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Function {
    return new Function().fromJsonString(jsonString, options);
  }

  static equals(a: Function | PlainMessage<Function> | undefined, b: Function | PlainMessage<Function> | undefined): boolean {
    return proto3.util.equals(Function, a, b);
  }
}

/**
 * @generated from message workflow.Input
 */
export class Input extends Message<Input> {
  /**
   * @generated from field: map<string, string> params = 1;
   */
  params: { [key: string]: string } = {};

  constructor(data?: PartialMessage<Input>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.Input";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "params", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Input {
    return new Input().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Input {
    return new Input().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Input {
    return new Input().fromJsonString(jsonString, options);
  }

  static equals(a: Input | PlainMessage<Input> | undefined, b: Input | PlainMessage<Input> | undefined): boolean {
    return proto3.util.equals(Input, a, b);
  }
}

/**
 * @generated from message workflow.Collection
 */
export class Collection extends Message<Collection> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  constructor(data?: PartialMessage<Collection>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.Collection";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Collection {
    return new Collection().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Collection {
    return new Collection().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Collection {
    return new Collection().fromJsonString(jsonString, options);
  }

  static equals(a: Collection | PlainMessage<Collection> | undefined, b: Collection | PlainMessage<Collection> | undefined): boolean {
    return proto3.util.equals(Collection, a, b);
  }
}

/**
 * @generated from message workflow.Table
 */
export class Table extends Message<Table> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  constructor(data?: PartialMessage<Table>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.Table";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Table {
    return new Table().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Table {
    return new Table().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Table {
    return new Table().fromJsonString(jsonString, options);
  }

  static equals(a: Table | PlainMessage<Table> | undefined, b: Table | PlainMessage<Table> | undefined): boolean {
    return proto3.util.equals(Table, a, b);
  }
}

/**
 * @generated from message workflow.Data
 */
export class Data extends Message<Data> {
  /**
   * @generated from oneof workflow.Data.type
   */
  type: {
    /**
     * @generated from field: workflow.Input input = 1;
     */
    value: Input;
    case: "input";
  } | {
    /**
     * @generated from field: workflow.Collection collection = 2;
     */
    value: Collection;
    case: "collection";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Data>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.Data";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "input", kind: "message", T: Input, oneof: "type" },
    { no: 2, name: "collection", kind: "message", T: Collection, oneof: "type" },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Data {
    return new Data().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Data {
    return new Data().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Data {
    return new Data().fromJsonString(jsonString, options);
  }

  static equals(a: Data | PlainMessage<Data> | undefined, b: Data | PlainMessage<Data> | undefined): boolean {
    return proto3.util.equals(Data, a, b);
  }
}

/**
 * @generated from message workflow.Node
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
   * @generated from oneof workflow.Node.type
   */
  type: {
    /**
     * @generated from field: workflow.Function function = 3;
     */
    value: Function;
    case: "function";
  } | {
    /**
     * @generated from field: workflow.Data data = 4;
     */
    value: Data;
    case: "data";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Node>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.Node";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "function", kind: "message", T: Function, oneof: "type" },
    { no: 4, name: "data", kind: "message", T: Data, oneof: "type" },
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
 * @generated from message workflow.Edge
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
  static readonly typeName = "workflow.Edge";
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

/**
 * @generated from message workflow.GRPCFunction
 */
export class GRPCFunction extends Message<GRPCFunction> {
  /**
   * @generated from field: string host = 1;
   */
  host = "";

  /**
   * @generated from field: string service = 2;
   */
  service = "";

  /**
   * @generated from field: string method = 3;
   */
  method = "";

  constructor(data?: PartialMessage<GRPCFunction>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.GRPCFunction";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "host", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "service", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "method", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GRPCFunction {
    return new GRPCFunction().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GRPCFunction {
    return new GRPCFunction().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GRPCFunction {
    return new GRPCFunction().fromJsonString(jsonString, options);
  }

  static equals(a: GRPCFunction | PlainMessage<GRPCFunction> | undefined, b: GRPCFunction | PlainMessage<GRPCFunction> | undefined): boolean {
    return proto3.util.equals(GRPCFunction, a, b);
  }
}

/**
 * @generated from message workflow.CodeFunction
 */
export class CodeFunction extends Message<CodeFunction> {
  /**
   * @generated from field: string code = 1;
   */
  code = "";

  constructor(data?: PartialMessage<CodeFunction>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.CodeFunction";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "code", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CodeFunction {
    return new CodeFunction().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CodeFunction {
    return new CodeFunction().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CodeFunction {
    return new CodeFunction().fromJsonString(jsonString, options);
  }

  static equals(a: CodeFunction | PlainMessage<CodeFunction> | undefined, b: CodeFunction | PlainMessage<CodeFunction> | undefined): boolean {
    return proto3.util.equals(CodeFunction, a, b);
  }
}

/**
 * @generated from message workflow.SQLFunction
 */
export class SQLFunction extends Message<SQLFunction> {
  /**
   * @generated from field: string url = 1;
   */
  url = "";

  /**
   * @generated from field: string query = 2;
   */
  query = "";

  constructor(data?: PartialMessage<SQLFunction>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.SQLFunction";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "query", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SQLFunction {
    return new SQLFunction().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SQLFunction {
    return new SQLFunction().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SQLFunction {
    return new SQLFunction().fromJsonString(jsonString, options);
  }

  static equals(a: SQLFunction | PlainMessage<SQLFunction> | undefined, b: SQLFunction | PlainMessage<SQLFunction> | undefined): boolean {
    return proto3.util.equals(SQLFunction, a, b);
  }
}

/**
 * @generated from message workflow.DBResource
 */
export class DBResource extends Message<DBResource> {
  /**
   * @generated from field: string url = 1;
   */
  url = "";

  constructor(data?: PartialMessage<DBResource>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.DBResource";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DBResource {
    return new DBResource().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DBResource {
    return new DBResource().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DBResource {
    return new DBResource().fromJsonString(jsonString, options);
  }

  static equals(a: DBResource | PlainMessage<DBResource> | undefined, b: DBResource | PlainMessage<DBResource> | undefined): boolean {
    return proto3.util.equals(DBResource, a, b);
  }
}

/**
 * @generated from message workflow.DocStoreResource
 */
export class DocStoreResource extends Message<DocStoreResource> {
  /**
   * @generated from field: string url = 1;
   */
  url = "";

  constructor(data?: PartialMessage<DocStoreResource>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.DocStoreResource";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DocStoreResource {
    return new DocStoreResource().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DocStoreResource {
    return new DocStoreResource().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DocStoreResource {
    return new DocStoreResource().fromJsonString(jsonString, options);
  }

  static equals(a: DocStoreResource | PlainMessage<DocStoreResource> | undefined, b: DocStoreResource | PlainMessage<DocStoreResource> | undefined): boolean {
    return proto3.util.equals(DocStoreResource, a, b);
  }
}

/**
 * @generated from message workflow.BucketResource
 */
export class BucketResource extends Message<BucketResource> {
  /**
   * @generated from field: string url = 1;
   */
  url = "";

  constructor(data?: PartialMessage<BucketResource>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "workflow.BucketResource";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): BucketResource {
    return new BucketResource().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): BucketResource {
    return new BucketResource().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): BucketResource {
    return new BucketResource().fromJsonString(jsonString, options);
  }

  static equals(a: BucketResource | PlainMessage<BucketResource> | undefined, b: BucketResource | PlainMessage<BucketResource> | undefined): boolean {
    return proto3.util.equals(BucketResource, a, b);
  }
}

