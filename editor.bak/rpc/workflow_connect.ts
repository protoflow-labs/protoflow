// @generated by protoc-gen-connect-es v0.8.6 with parameter "target=ts"
// @generated from file workflow.proto (package workflow, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { ID, Run, Workflow, WorkflowEntrypoint } from "./workflow_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service workflow.ManagerService
 */
export const ManagerService = {
  typeName: "workflow.ManagerService",
  methods: {
    /**
     * @generated from rpc workflow.ManagerService.CreateWorkflow
     */
    createWorkflow: {
      name: "CreateWorkflow",
      I: Workflow,
      O: ID,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc workflow.ManagerService.StartWorkflow
     */
    startWorkflow: {
      name: "StartWorkflow",
      I: WorkflowEntrypoint,
      O: Run,
      kind: MethodKind.Unary,
    },
  }
} as const;
