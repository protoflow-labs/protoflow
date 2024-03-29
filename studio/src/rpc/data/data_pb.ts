// @generated by protoc-gen-es v1.5.1 with parameter "target=ts"
// @generated from file data/data.proto (package data, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message data.Input
 */
export class Input extends Message<Input> {
  /**
   * @generated from field: string value = 1;
   */
  value = "";

  constructor(data?: PartialMessage<Input>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "data.Input";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "value", kind: "scalar", T: 9 /* ScalarType.STRING */ },
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
 * TODO breadchris do we need an explicit provider for config? probably
 *
 * @generated from message data.Config
 */
export class Config extends Message<Config> {
  /**
   * TODO breadchris configs could have a type
   * google.protobuf.DescriptorProto type = 1;
   *
   * @generated from field: string value = 1;
   */
  value = "";

  constructor(data?: PartialMessage<Config>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "data.Config";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "value", kind: "scalar", T: 9 /* ScalarType.STRING */ },
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
 * @generated from message data.Data
 */
export class Data extends Message<Data> {
  /**
   * @generated from oneof data.Data.type
   */
  type: {
    /**
     * @generated from field: data.Input input = 1;
     */
    value: Input;
    case: "input";
  } | {
    /**
     * @generated from field: data.Config config = 2;
     */
    value: Config;
    case: "config";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Data>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "data.Data";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "input", kind: "message", T: Input, oneof: "type" },
    { no: 2, name: "config", kind: "message", T: Config, oneof: "type" },
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

