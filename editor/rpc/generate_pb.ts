// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file generate.proto (package generate, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message generate.GenerateRequest
 */
export class GenerateRequest extends Message<GenerateRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  constructor(data?: PartialMessage<GenerateRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "generate.GenerateRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GenerateRequest {
    return new GenerateRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GenerateRequest {
    return new GenerateRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GenerateRequest {
    return new GenerateRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GenerateRequest | PlainMessage<GenerateRequest> | undefined, b: GenerateRequest | PlainMessage<GenerateRequest> | undefined): boolean {
    return proto3.util.equals(GenerateRequest, a, b);
  }
}

/**
 * @generated from message generate.GenerateResponse
 */
export class GenerateResponse extends Message<GenerateResponse> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  constructor(data?: PartialMessage<GenerateResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "generate.GenerateResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GenerateResponse {
    return new GenerateResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GenerateResponse {
    return new GenerateResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GenerateResponse {
    return new GenerateResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GenerateResponse | PlainMessage<GenerateResponse> | undefined, b: GenerateResponse | PlainMessage<GenerateResponse> | undefined): boolean {
    return proto3.util.equals(GenerateResponse, a, b);
  }
}

