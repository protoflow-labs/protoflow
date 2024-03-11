// @generated by protoc-gen-es v1.5.1 with parameter "target=ts"
// @generated from file ai.proto (package ai, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message ai.GenerateCode
 */
export class GenerateCode extends Message<GenerateCode> {
  /**
   * Code that has been generated.
   *
   * @generated from field: string code = 1;
   */
  code = "";

  constructor(data?: PartialMessage<GenerateCode>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "ai.GenerateCode";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "code", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GenerateCode {
    return new GenerateCode().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GenerateCode {
    return new GenerateCode().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GenerateCode {
    return new GenerateCode().fromJsonString(jsonString, options);
  }

  static equals(a: GenerateCode | PlainMessage<GenerateCode> | undefined, b: GenerateCode | PlainMessage<GenerateCode> | undefined): boolean {
    return proto3.util.equals(GenerateCode, a, b);
  }
}

/**
 * @generated from message ai.DirectionResponse
 */
export class DirectionResponse extends Message<DirectionResponse> {
  /**
   * @generated from field: repeated string ingredients = 1;
   */
  ingredients: string[] = [];

  /**
   * @generated from field: repeated string quantities = 2;
   */
  quantities: string[] = [];

  /**
   * @generated from field: repeated string steps = 3;
   */
  steps: string[] = [];

  constructor(data?: PartialMessage<DirectionResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "ai.DirectionResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "ingredients", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 2, name: "quantities", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 3, name: "steps", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DirectionResponse {
    return new DirectionResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DirectionResponse {
    return new DirectionResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DirectionResponse {
    return new DirectionResponse().fromJsonString(jsonString, options);
  }

  static equals(a: DirectionResponse | PlainMessage<DirectionResponse> | undefined, b: DirectionResponse | PlainMessage<DirectionResponse> | undefined): boolean {
    return proto3.util.equals(DirectionResponse, a, b);
  }
}

