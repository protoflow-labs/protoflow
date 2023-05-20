// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file block.proto (package block, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { DescriptorProto, EnumDescriptorProto, Message, MethodDescriptorProto, proto3 } from "@bufbuild/protobuf";

/**
 * TODO breadchris think through this more
 *
 * @generated from message block.Block
 */
export class Block extends Message<Block> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: string version = 3;
   */
  version = "";

  /**
   * @generated from oneof block.Block.type
   */
  type: {
    /**
     * @generated from field: block.REST rest = 4;
     */
    value: REST;
    case: "rest";
  } | {
    /**
     * @generated from field: block.GRPC grpc = 5;
     */
    value: GRPC;
    case: "grpc";
  } | {
    /**
     * @generated from field: block.Collection collection = 6;
     */
    value: Collection;
    case: "collection";
  } | {
    /**
     * @generated from field: block.Input input = 7;
     */
    value: Input;
    case: "input";
  } | {
    /**
     * @generated from field: block.Bucket bucket = 8;
     */
    value: Bucket;
    case: "bucket";
  } | {
    /**
     * @generated from field: block.Function function = 9;
     */
    value: Function;
    case: "function";
  } | {
    /**
     * @generated from field: block.Query query = 10;
     */
    value: Query;
    case: "query";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Block>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "block.Block";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "version", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "rest", kind: "message", T: REST, oneof: "type" },
    { no: 5, name: "grpc", kind: "message", T: GRPC, oneof: "type" },
    { no: 6, name: "collection", kind: "message", T: Collection, oneof: "type" },
    { no: 7, name: "input", kind: "message", T: Input, oneof: "type" },
    { no: 8, name: "bucket", kind: "message", T: Bucket, oneof: "type" },
    { no: 9, name: "function", kind: "message", T: Function, oneof: "type" },
    { no: 10, name: "query", kind: "message", T: Query, oneof: "type" },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Block {
    return new Block().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Block {
    return new Block().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Block {
    return new Block().fromJsonString(jsonString, options);
  }

  static equals(a: Block | PlainMessage<Block> | undefined, b: Block | PlainMessage<Block> | undefined): boolean {
    return proto3.util.equals(Block, a, b);
  }
}

/**
 * @generated from message block.Input
 */
export class Input extends Message<Input> {
  /**
   * @generated from field: repeated block.FieldDefinition fields = 1;
   */
  fields: FieldDefinition[] = [];

  constructor(data?: PartialMessage<Input>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "block.Input";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "fields", kind: "message", T: FieldDefinition, repeated: true },
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
 * @generated from message block.Config
 */
export class Config extends Message<Config> {
  /**
   * @generated from field: google.protobuf.DescriptorProto type = 1;
   */
  type?: DescriptorProto;

  constructor(data?: PartialMessage<Config>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "block.Config";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "type", kind: "message", T: DescriptorProto },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Config {
    return new Config().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Config {
    return new Config().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Config {
    return new Config().fromJsonString(jsonString, options);
  }

  static equals(a: Config | PlainMessage<Config> | undefined, b: Config | PlainMessage<Config> | undefined): boolean {
    return proto3.util.equals(Config, a, b);
  }
}

/**
 * @generated from message block.Collection
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
  static readonly typeName = "block.Collection";
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
 * @generated from message block.Bucket
 */
export class Bucket extends Message<Bucket> {
  /**
   * @generated from field: string path = 1;
   */
  path = "";

  constructor(data?: PartialMessage<Bucket>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "block.Bucket";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "path", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Bucket {
    return new Bucket().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Bucket {
    return new Bucket().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Bucket {
    return new Bucket().fromJsonString(jsonString, options);
  }

  static equals(a: Bucket | PlainMessage<Bucket> | undefined, b: Bucket | PlainMessage<Bucket> | undefined): boolean {
    return proto3.util.equals(Bucket, a, b);
  }
}

/**
 * @generated from message block.Function
 */
export class Function extends Message<Function> {
  /**
   * @generated from field: string runtime = 1;
   */
  runtime = "";

  /**
   * @generated from field: block.GRPC grpc = 2;
   */
  grpc?: GRPC;

  constructor(data?: PartialMessage<Function>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "block.Function";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "runtime", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "grpc", kind: "message", T: GRPC },
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
 * @generated from message block.Query
 */
export class Query extends Message<Query> {
  /**
   * @generated from field: string collection = 1;
   */
  collection = "";

  constructor(data?: PartialMessage<Query>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "block.Query";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "collection", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Query {
    return new Query().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Query {
    return new Query().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Query {
    return new Query().fromJsonString(jsonString, options);
  }

  static equals(a: Query | PlainMessage<Query> | undefined, b: Query | PlainMessage<Query> | undefined): boolean {
    return proto3.util.equals(Query, a, b);
  }
}

/**
 * @generated from message block.Result
 */
export class Result extends Message<Result> {
  /**
   * @generated from field: bytes data = 1;
   */
  data = new Uint8Array(0);

  constructor(data?: PartialMessage<Result>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "block.Result";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "data", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Result {
    return new Result().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Result {
    return new Result().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Result {
    return new Result().fromJsonString(jsonString, options);
  }

  static equals(a: Result | PlainMessage<Result> | undefined, b: Result | PlainMessage<Result> | undefined): boolean {
    return proto3.util.equals(Result, a, b);
  }
}

/**
 * @generated from message block.FieldDefinition
 */
export class FieldDefinition extends Message<FieldDefinition> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: block.FieldDefinition.FieldType type = 2;
   */
  type = FieldDefinition_FieldType.STRING;

  constructor(data?: PartialMessage<FieldDefinition>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "block.FieldDefinition";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "type", kind: "enum", T: proto3.getEnumType(FieldDefinition_FieldType) },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): FieldDefinition {
    return new FieldDefinition().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): FieldDefinition {
    return new FieldDefinition().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): FieldDefinition {
    return new FieldDefinition().fromJsonString(jsonString, options);
  }

  static equals(a: FieldDefinition | PlainMessage<FieldDefinition> | undefined, b: FieldDefinition | PlainMessage<FieldDefinition> | undefined): boolean {
    return proto3.util.equals(FieldDefinition, a, b);
  }
}

/**
 * @generated from enum block.FieldDefinition.FieldType
 */
export enum FieldDefinition_FieldType {
  /**
   * @generated from enum value: STRING = 0;
   */
  STRING = 0,

  /**
   * @generated from enum value: INTEGER = 1;
   */
  INTEGER = 1,

  /**
   * @generated from enum value: BOOLEAN = 2;
   */
  BOOLEAN = 2,
}
// Retrieve enum metadata with: proto3.getEnumType(FieldDefinition_FieldType)
proto3.util.setEnumType(FieldDefinition_FieldType, "block.FieldDefinition.FieldType", [
  { no: 0, name: "STRING" },
  { no: 1, name: "INTEGER" },
  { no: 2, name: "BOOLEAN" },
]);

/**
 * @generated from message block.REST
 */
export class REST extends Message<REST> {
  /**
   * @generated from field: string path = 1;
   */
  path = "";

  /**
   * @generated from field: string method = 2;
   */
  method = "";

  /**
   * @generated from field: map<string, string> headers = 3;
   */
  headers: { [key: string]: string } = {};

  constructor(data?: PartialMessage<REST>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "block.REST";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "path", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "method", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "headers", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): REST {
    return new REST().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): REST {
    return new REST().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): REST {
    return new REST().fromJsonString(jsonString, options);
  }

  static equals(a: REST | PlainMessage<REST> | undefined, b: REST | PlainMessage<REST> | undefined): boolean {
    return proto3.util.equals(REST, a, b);
  }
}

/**
 * @generated from message block.GRPC
 */
export class GRPC extends Message<GRPC> {
  /**
   * @generated from field: string package = 1;
   */
  package = "";

  /**
   * @generated from field: string service = 2;
   */
  service = "";

  /**
   * @generated from field: string method = 3;
   */
  method = "";

  /**
   * @generated from field: google.protobuf.DescriptorProto input = 4;
   */
  input?: DescriptorProto;

  /**
   * @generated from field: google.protobuf.DescriptorProto output = 5;
   */
  output?: DescriptorProto;

  /**
   * @generated from field: map<string, google.protobuf.DescriptorProto> desc_lookup = 6;
   */
  descLookup: { [key: string]: DescriptorProto } = {};

  /**
   * @generated from field: map<string, google.protobuf.EnumDescriptorProto> enum_lookup = 7;
   */
  enumLookup: { [key: string]: EnumDescriptorProto } = {};

  /**
   * @generated from field: google.protobuf.MethodDescriptorProto method_desc = 8;
   */
  methodDesc?: MethodDescriptorProto;

  constructor(data?: PartialMessage<GRPC>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "block.GRPC";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "package", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "service", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "method", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "input", kind: "message", T: DescriptorProto },
    { no: 5, name: "output", kind: "message", T: DescriptorProto },
    { no: 6, name: "desc_lookup", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "message", T: DescriptorProto} },
    { no: 7, name: "enum_lookup", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "message", T: EnumDescriptorProto} },
    { no: 8, name: "method_desc", kind: "message", T: MethodDescriptorProto },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GRPC {
    return new GRPC().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GRPC {
    return new GRPC().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GRPC {
    return new GRPC().fromJsonString(jsonString, options);
  }

  static equals(a: GRPC | PlainMessage<GRPC> | undefined, b: GRPC | PlainMessage<GRPC> | undefined): boolean {
    return proto3.util.equals(GRPC, a, b);
  }
}

