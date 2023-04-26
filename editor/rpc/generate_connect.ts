// @generated by protoc-gen-connect-es v0.8.6 with parameter "target=ts"
// @generated from file generate.proto (package generate, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { GenerateRequest, GenerateResponse } from "./generate_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service generate.GenerateService
 */
export const GenerateService = {
  typeName: "generate.GenerateService",
  methods: {
    /**
     * @generated from rpc generate.GenerateService.Generate
     */
    generate: {
      name: "Generate",
      I: GenerateRequest,
      O: GenerateResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;
