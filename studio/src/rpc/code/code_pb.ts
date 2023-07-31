// @generated by protoc-gen-es v1.2.1 with parameter "target=ts"
// @generated from file code/code.proto (package code, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { Server as Server$1 } from "../grpc/grpc_pb.js";

/**
 * @generated from enum code.Runtime
 */
export enum Runtime {
  /**
   * @generated from enum value: NODEJS = 0;
   */
  NODEJS = 0,

  /**
   * @generated from enum value: PYTHON = 1;
   */
  PYTHON = 1,

  /**
   * @generated from enum value: GO = 2;
   */
  GO = 2,
}
// Retrieve enum metadata with: proto3.getEnumType(Runtime)
proto3.util.setEnumType(Runtime, "code.Runtime", [
  { no: 0, name: "NODEJS" },
  { no: 1, name: "PYTHON" },
  { no: 2, name: "GO" },
]);

/**
 * @generated from message code.Function
 */
export class Function extends Message<Function> {
  constructor(data?: PartialMessage<Function>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "code.Function";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
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
 * @generated from message code.Server
 */
export class Server extends Message<Server> {
  /**
   * @generated from field: code.Runtime runtime = 1;
   */
  runtime = Runtime.NODEJS;

  /**
   * string containerURI = 4;
   *
   * @generated from field: grpc.Server grpc = 2;
   */
  grpc?: Server$1;

  constructor(data?: PartialMessage<Server>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "code.Server";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "runtime", kind: "enum", T: proto3.getEnumType(Runtime) },
    { no: 2, name: "grpc", kind: "message", T: Server$1 },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Server {
    return new Server().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Server {
    return new Server().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Server {
    return new Server().fromJsonString(jsonString, options);
  }

  static equals(a: Server | PlainMessage<Server> | undefined, b: Server | PlainMessage<Server> | undefined): boolean {
    return proto3.util.equals(Server, a, b);
  }
}

/**
 * @generated from message code.Code
 */
export class Code extends Message<Code> {
  /**
   * @generated from oneof code.Code.type
   */
  type: {
    /**
     * @generated from field: code.Function function = 1;
     */
    value: Function;
    case: "function";
  } | {
    /**
     * @generated from field: code.Server server = 2;
     */
    value: Server;
    case: "server";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Code>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "code.Code";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "function", kind: "message", T: Function, oneof: "type" },
    { no: 2, name: "server", kind: "message", T: Server, oneof: "type" },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Code {
    return new Code().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Code {
    return new Code().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Code {
    return new Code().fromJsonString(jsonString, options);
  }

  static equals(a: Code | PlainMessage<Code> | undefined, b: Code | PlainMessage<Code> | undefined): boolean {
    return proto3.util.equals(Code, a, b);
  }
}

