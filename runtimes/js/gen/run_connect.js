// @generated by protoc-gen-connect-es v0.8.6 with parameter "target=js"
// @generated from file run.proto (package run, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { RunRequest, RunResponse } from "./run_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service run.RunService
 */
export const RunService = {
  typeName: "run.RunService",
  methods: {
    /**
     * @generated from rpc run.RunService.Run
     */
    run: {
      name: "Run",
      I: RunRequest,
      O: RunResponse,
      kind: MethodKind.Unary,
    },
  }
};

