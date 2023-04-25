// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file project.proto (package project, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Any, Message, proto3 } from "@bufbuild/protobuf";
import { Graph } from "./graph_pb.js";
import { Block } from "./block_pb.js";
import { Resource } from "./resource_pb.js";

/**
 * @generated from message project.RunWorkflowRequest
 */
export class RunWorkflowRequest extends Message<RunWorkflowRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  /**
   * @generated from field: string node_id = 2;
   */
  nodeId = "";

  /**
   * @generated from field: google.protobuf.Any input = 3;
   */
  input?: Any;

  constructor(data?: PartialMessage<RunWorkflowRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.RunWorkflowRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "node_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "input", kind: "message", T: Any },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RunWorkflowRequest {
    return new RunWorkflowRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RunWorkflowRequest {
    return new RunWorkflowRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RunWorkflowRequest {
    return new RunWorkflowRequest().fromJsonString(jsonString, options);
  }

  static equals(a: RunWorkflowRequest | PlainMessage<RunWorkflowRequest> | undefined, b: RunWorkflowRequest | PlainMessage<RunWorkflowRequest> | undefined): boolean {
    return proto3.util.equals(RunWorkflowRequest, a, b);
  }
}

/**
 * @generated from message project.RunBlockRequest
 */
export class RunBlockRequest extends Message<RunBlockRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  /**
   * @generated from field: string block_id = 2;
   */
  blockId = "";

  /**
   * @generated from field: google.protobuf.Any input = 3;
   */
  input?: Any;

  constructor(data?: PartialMessage<RunBlockRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.RunBlockRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "block_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "input", kind: "message", T: Any },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RunBlockRequest {
    return new RunBlockRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RunBlockRequest {
    return new RunBlockRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RunBlockRequest {
    return new RunBlockRequest().fromJsonString(jsonString, options);
  }

  static equals(a: RunBlockRequest | PlainMessage<RunBlockRequest> | undefined, b: RunBlockRequest | PlainMessage<RunBlockRequest> | undefined): boolean {
    return proto3.util.equals(RunBlockRequest, a, b);
  }
}

/**
 * @generated from message project.RunOutput
 */
export class RunOutput extends Message<RunOutput> {
  /**
   * @generated from field: google.protobuf.Any output = 1;
   */
  output?: Any;

  constructor(data?: PartialMessage<RunOutput>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.RunOutput";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "output", kind: "message", T: Any },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RunOutput {
    return new RunOutput().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RunOutput {
    return new RunOutput().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RunOutput {
    return new RunOutput().fromJsonString(jsonString, options);
  }

  static equals(a: RunOutput | PlainMessage<RunOutput> | undefined, b: RunOutput | PlainMessage<RunOutput> | undefined): boolean {
    return proto3.util.equals(RunOutput, a, b);
  }
}

/**
 * @generated from message project.Project
 */
export class Project extends Message<Project> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: string description = 3;
   */
  description = "";

  /**
   * @generated from field: string owner = 4;
   */
  owner = "";

  /**
   * @generated from field: string created_at = 5;
   */
  createdAt = "";

  /**
   * @generated from field: string updated_at = 6;
   */
  updatedAt = "";

  /**
   * @generated from field: graph.Graph graph = 7;
   */
  graph?: Graph;

  /**
   * @generated from field: repeated block.Block blocks = 8;
   */
  blocks: Block[] = [];

  /**
   * @generated from field: repeated resource.Resource resources = 9;
   */
  resources: Resource[] = [];

  constructor(data?: PartialMessage<Project>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.Project";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "owner", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "created_at", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "updated_at", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 7, name: "graph", kind: "message", T: Graph },
    { no: 8, name: "blocks", kind: "message", T: Block, repeated: true },
    { no: 9, name: "resources", kind: "message", T: Resource, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Project {
    return new Project().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Project {
    return new Project().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Project {
    return new Project().fromJsonString(jsonString, options);
  }

  static equals(a: Project | PlainMessage<Project> | undefined, b: Project | PlainMessage<Project> | undefined): boolean {
    return proto3.util.equals(Project, a, b);
  }
}

/**
 * @generated from message project.CreateResourceRequest
 */
export class CreateResourceRequest extends Message<CreateResourceRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  /**
   * @generated from field: resource.Resource resource = 2;
   */
  resource?: Resource;

  constructor(data?: PartialMessage<CreateResourceRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.CreateResourceRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "resource", kind: "message", T: Resource },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateResourceRequest {
    return new CreateResourceRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateResourceRequest {
    return new CreateResourceRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateResourceRequest {
    return new CreateResourceRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CreateResourceRequest | PlainMessage<CreateResourceRequest> | undefined, b: CreateResourceRequest | PlainMessage<CreateResourceRequest> | undefined): boolean {
    return proto3.util.equals(CreateResourceRequest, a, b);
  }
}

/**
 * @generated from message project.CreateResourceResponse
 */
export class CreateResourceResponse extends Message<CreateResourceResponse> {
  /**
   * @generated from field: string resource_id = 1;
   */
  resourceId = "";

  constructor(data?: PartialMessage<CreateResourceResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.CreateResourceResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "resource_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateResourceResponse {
    return new CreateResourceResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateResourceResponse {
    return new CreateResourceResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateResourceResponse {
    return new CreateResourceResponse().fromJsonString(jsonString, options);
  }

  static equals(a: CreateResourceResponse | PlainMessage<CreateResourceResponse> | undefined, b: CreateResourceResponse | PlainMessage<CreateResourceResponse> | undefined): boolean {
    return proto3.util.equals(CreateResourceResponse, a, b);
  }
}

/**
 * @generated from message project.GetProjectRequest
 */
export class GetProjectRequest extends Message<GetProjectRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  constructor(data?: PartialMessage<GetProjectRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.GetProjectRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetProjectRequest {
    return new GetProjectRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetProjectRequest {
    return new GetProjectRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetProjectRequest {
    return new GetProjectRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetProjectRequest | PlainMessage<GetProjectRequest> | undefined, b: GetProjectRequest | PlainMessage<GetProjectRequest> | undefined): boolean {
    return proto3.util.equals(GetProjectRequest, a, b);
  }
}

/**
 * @generated from message project.GetProjectResponse
 */
export class GetProjectResponse extends Message<GetProjectResponse> {
  /**
   * @generated from field: project.Project project = 1;
   */
  project?: Project;

  constructor(data?: PartialMessage<GetProjectResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.GetProjectResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project", kind: "message", T: Project },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetProjectResponse {
    return new GetProjectResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetProjectResponse {
    return new GetProjectResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetProjectResponse {
    return new GetProjectResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetProjectResponse | PlainMessage<GetProjectResponse> | undefined, b: GetProjectResponse | PlainMessage<GetProjectResponse> | undefined): boolean {
    return proto3.util.equals(GetProjectResponse, a, b);
  }
}

/**
 * @generated from message project.GetProjectsRequest
 */
export class GetProjectsRequest extends Message<GetProjectsRequest> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  constructor(data?: PartialMessage<GetProjectsRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.GetProjectsRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetProjectsRequest {
    return new GetProjectsRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetProjectsRequest {
    return new GetProjectsRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetProjectsRequest {
    return new GetProjectsRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetProjectsRequest | PlainMessage<GetProjectsRequest> | undefined, b: GetProjectsRequest | PlainMessage<GetProjectsRequest> | undefined): boolean {
    return proto3.util.equals(GetProjectsRequest, a, b);
  }
}

/**
 * @generated from message project.GetProjectsResponse
 */
export class GetProjectsResponse extends Message<GetProjectsResponse> {
  /**
   * @generated from field: repeated project.Project projects = 1;
   */
  projects: Project[] = [];

  constructor(data?: PartialMessage<GetProjectsResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.GetProjectsResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "projects", kind: "message", T: Project, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetProjectsResponse {
    return new GetProjectsResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetProjectsResponse {
    return new GetProjectsResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetProjectsResponse {
    return new GetProjectsResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetProjectsResponse | PlainMessage<GetProjectsResponse> | undefined, b: GetProjectsResponse | PlainMessage<GetProjectsResponse> | undefined): boolean {
    return proto3.util.equals(GetProjectsResponse, a, b);
  }
}

/**
 * @generated from message project.CreateProjectRequest
 */
export class CreateProjectRequest extends Message<CreateProjectRequest> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  constructor(data?: PartialMessage<CreateProjectRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.CreateProjectRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateProjectRequest {
    return new CreateProjectRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateProjectRequest {
    return new CreateProjectRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateProjectRequest {
    return new CreateProjectRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CreateProjectRequest | PlainMessage<CreateProjectRequest> | undefined, b: CreateProjectRequest | PlainMessage<CreateProjectRequest> | undefined): boolean {
    return proto3.util.equals(CreateProjectRequest, a, b);
  }
}

/**
 * @generated from message project.CreateProjectResponse
 */
export class CreateProjectResponse extends Message<CreateProjectResponse> {
  /**
   * @generated from field: project.Project project = 1;
   */
  project?: Project;

  constructor(data?: PartialMessage<CreateProjectResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.CreateProjectResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project", kind: "message", T: Project },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateProjectResponse {
    return new CreateProjectResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateProjectResponse {
    return new CreateProjectResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateProjectResponse {
    return new CreateProjectResponse().fromJsonString(jsonString, options);
  }

  static equals(a: CreateProjectResponse | PlainMessage<CreateProjectResponse> | undefined, b: CreateProjectResponse | PlainMessage<CreateProjectResponse> | undefined): boolean {
    return proto3.util.equals(CreateProjectResponse, a, b);
  }
}

/**
 * @generated from message project.DeleteProjectRequest
 */
export class DeleteProjectRequest extends Message<DeleteProjectRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  constructor(data?: PartialMessage<DeleteProjectRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.DeleteProjectRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeleteProjectRequest {
    return new DeleteProjectRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeleteProjectRequest {
    return new DeleteProjectRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeleteProjectRequest {
    return new DeleteProjectRequest().fromJsonString(jsonString, options);
  }

  static equals(a: DeleteProjectRequest | PlainMessage<DeleteProjectRequest> | undefined, b: DeleteProjectRequest | PlainMessage<DeleteProjectRequest> | undefined): boolean {
    return proto3.util.equals(DeleteProjectRequest, a, b);
  }
}

/**
 * @generated from message project.DeleteProjectResponse
 */
export class DeleteProjectResponse extends Message<DeleteProjectResponse> {
  /**
   * @generated from field: project.Project project = 1;
   */
  project?: Project;

  constructor(data?: PartialMessage<DeleteProjectResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.DeleteProjectResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project", kind: "message", T: Project },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeleteProjectResponse {
    return new DeleteProjectResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeleteProjectResponse {
    return new DeleteProjectResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeleteProjectResponse {
    return new DeleteProjectResponse().fromJsonString(jsonString, options);
  }

  static equals(a: DeleteProjectResponse | PlainMessage<DeleteProjectResponse> | undefined, b: DeleteProjectResponse | PlainMessage<DeleteProjectResponse> | undefined): boolean {
    return proto3.util.equals(DeleteProjectResponse, a, b);
  }
}

/**
 * @generated from message project.GetResourcesRequest
 */
export class GetResourcesRequest extends Message<GetResourcesRequest> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  constructor(data?: PartialMessage<GetResourcesRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.GetResourcesRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetResourcesRequest {
    return new GetResourcesRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetResourcesRequest {
    return new GetResourcesRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetResourcesRequest {
    return new GetResourcesRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetResourcesRequest | PlainMessage<GetResourcesRequest> | undefined, b: GetResourcesRequest | PlainMessage<GetResourcesRequest> | undefined): boolean {
    return proto3.util.equals(GetResourcesRequest, a, b);
  }
}

/**
 * @generated from message project.GetResourcesResponse
 */
export class GetResourcesResponse extends Message<GetResourcesResponse> {
  /**
   * @generated from field: repeated resource.Resource resources = 1;
   */
  resources: Resource[] = [];

  constructor(data?: PartialMessage<GetResourcesResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.GetResourcesResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "resources", kind: "message", T: Resource, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetResourcesResponse {
    return new GetResourcesResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetResourcesResponse {
    return new GetResourcesResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetResourcesResponse {
    return new GetResourcesResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetResourcesResponse | PlainMessage<GetResourcesResponse> | undefined, b: GetResourcesResponse | PlainMessage<GetResourcesResponse> | undefined): boolean {
    return proto3.util.equals(GetResourcesResponse, a, b);
  }
}

/**
 * @generated from message project.SaveProjectRequest
 */
export class SaveProjectRequest extends Message<SaveProjectRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  /**
   * @generated from field: graph.Graph graph = 2;
   */
  graph?: Graph;

  /**
   * @generated from field: repeated block.Block blocks = 3;
   */
  blocks: Block[] = [];

  /**
   * @generated from field: repeated resource.Resource resources = 4;
   */
  resources: Resource[] = [];

  constructor(data?: PartialMessage<SaveProjectRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.SaveProjectRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "graph", kind: "message", T: Graph },
    { no: 3, name: "blocks", kind: "message", T: Block, repeated: true },
    { no: 4, name: "resources", kind: "message", T: Resource, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SaveProjectRequest {
    return new SaveProjectRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SaveProjectRequest {
    return new SaveProjectRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SaveProjectRequest {
    return new SaveProjectRequest().fromJsonString(jsonString, options);
  }

  static equals(a: SaveProjectRequest | PlainMessage<SaveProjectRequest> | undefined, b: SaveProjectRequest | PlainMessage<SaveProjectRequest> | undefined): boolean {
    return proto3.util.equals(SaveProjectRequest, a, b);
  }
}

/**
 * @generated from message project.SaveProjectResponse
 */
export class SaveProjectResponse extends Message<SaveProjectResponse> {
  /**
   * @generated from field: project.Project project = 1;
   */
  project?: Project;

  constructor(data?: PartialMessage<SaveProjectResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.SaveProjectResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project", kind: "message", T: Project },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SaveProjectResponse {
    return new SaveProjectResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SaveProjectResponse {
    return new SaveProjectResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SaveProjectResponse {
    return new SaveProjectResponse().fromJsonString(jsonString, options);
  }

  static equals(a: SaveProjectResponse | PlainMessage<SaveProjectResponse> | undefined, b: SaveProjectResponse | PlainMessage<SaveProjectResponse> | undefined): boolean {
    return proto3.util.equals(SaveProjectResponse, a, b);
  }
}

