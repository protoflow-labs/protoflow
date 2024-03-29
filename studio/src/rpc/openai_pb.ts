// @generated by protoc-gen-es v1.5.1 with parameter "target=ts"
// @generated from file openai.proto (package openai, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message openai.PromptRequest
 */
export class PromptRequest extends Message<PromptRequest> {
  /**
   * @generated from field: string prompt = 1;
   */
  prompt = "";

  constructor(data?: PartialMessage<PromptRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "openai.PromptRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "prompt", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PromptRequest {
    return new PromptRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PromptRequest {
    return new PromptRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PromptRequest {
    return new PromptRequest().fromJsonString(jsonString, options);
  }

  static equals(a: PromptRequest | PlainMessage<PromptRequest> | undefined, b: PromptRequest | PlainMessage<PromptRequest> | undefined): boolean {
    return proto3.util.equals(PromptRequest, a, b);
  }
}

/**
 * @generated from message openai.PromptResponse
 */
export class PromptResponse extends Message<PromptResponse> {
  /**
   * @generated from field: string text = 1;
   */
  text = "";

  constructor(data?: PartialMessage<PromptResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "openai.PromptResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "text", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PromptResponse {
    return new PromptResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PromptResponse {
    return new PromptResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PromptResponse {
    return new PromptResponse().fromJsonString(jsonString, options);
  }

  static equals(a: PromptResponse | PlainMessage<PromptResponse> | undefined, b: PromptResponse | PlainMessage<PromptResponse> | undefined): boolean {
    return proto3.util.equals(PromptResponse, a, b);
  }
}

