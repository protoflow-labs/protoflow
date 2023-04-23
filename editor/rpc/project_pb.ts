// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file project.proto (package project, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message project.Resource
 */
export class Resource extends Message<Resource> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: string image = 3;
   */
  image = "";

  constructor(data?: PartialMessage<Resource>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.Resource";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "image", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Resource {
    return new Resource().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Resource {
    return new Resource().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Resource {
    return new Resource().fromJsonString(jsonString, options);
  }

  static equals(a: Resource | PlainMessage<Resource> | undefined, b: Resource | PlainMessage<Resource> | undefined): boolean {
    return proto3.util.equals(Resource, a, b);
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
 * @generated from message project.Block
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
   * @generated from field: string type = 3;
   */
  type = "";

  /**
   * @generated from field: float x = 4;
   */
  x = 0;

  /**
   * @generated from field: float y = 5;
   */
  y = 0;

  constructor(data?: PartialMessage<Block>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.Block";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "type", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "x", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
    { no: 5, name: "y", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
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
 * @generated from message project.Edge
 */
export class Edge extends Message<Edge> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string source = 2;
   */
  source = "";

  /**
   * @generated from field: string target = 3;
   */
  target = "";

  constructor(data?: PartialMessage<Edge>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.Edge";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "source", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "target", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Edge {
    return new Edge().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Edge {
    return new Edge().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Edge {
    return new Edge().fromJsonString(jsonString, options);
  }

  static equals(a: Edge | PlainMessage<Edge> | undefined, b: Edge | PlainMessage<Edge> | undefined): boolean {
    return proto3.util.equals(Edge, a, b);
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
   * @generated from field: project.Project project = 1;
   */
  project?: Project;

  constructor(data?: PartialMessage<CreateProjectRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.CreateProjectRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project", kind: "message", T: Project },
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
   * @generated from field: repeated project.Resource resources = 1;
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
 * @generated from message project.AddBlockRequest
 */
export class AddBlockRequest extends Message<AddBlockRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  /**
   * @generated from field: project.Block block = 2;
   */
  block?: Block;

  constructor(data?: PartialMessage<AddBlockRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.AddBlockRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "block", kind: "message", T: Block },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AddBlockRequest {
    return new AddBlockRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AddBlockRequest {
    return new AddBlockRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AddBlockRequest {
    return new AddBlockRequest().fromJsonString(jsonString, options);
  }

  static equals(a: AddBlockRequest | PlainMessage<AddBlockRequest> | undefined, b: AddBlockRequest | PlainMessage<AddBlockRequest> | undefined): boolean {
    return proto3.util.equals(AddBlockRequest, a, b);
  }
}

/**
 * @generated from message project.AddBlockResponse
 */
export class AddBlockResponse extends Message<AddBlockResponse> {
  /**
   * @generated from field: project.Block block = 1;
   */
  block?: Block;

  constructor(data?: PartialMessage<AddBlockResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.AddBlockResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "block", kind: "message", T: Block },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AddBlockResponse {
    return new AddBlockResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AddBlockResponse {
    return new AddBlockResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AddBlockResponse {
    return new AddBlockResponse().fromJsonString(jsonString, options);
  }

  static equals(a: AddBlockResponse | PlainMessage<AddBlockResponse> | undefined, b: AddBlockResponse | PlainMessage<AddBlockResponse> | undefined): boolean {
    return proto3.util.equals(AddBlockResponse, a, b);
  }
}

/**
 * @generated from message project.RemoveBlockRequest
 */
export class RemoveBlockRequest extends Message<RemoveBlockRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  /**
   * @generated from field: string block_id = 2;
   */
  blockId = "";

  constructor(data?: PartialMessage<RemoveBlockRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.RemoveBlockRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "block_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RemoveBlockRequest {
    return new RemoveBlockRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RemoveBlockRequest {
    return new RemoveBlockRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RemoveBlockRequest {
    return new RemoveBlockRequest().fromJsonString(jsonString, options);
  }

  static equals(a: RemoveBlockRequest | PlainMessage<RemoveBlockRequest> | undefined, b: RemoveBlockRequest | PlainMessage<RemoveBlockRequest> | undefined): boolean {
    return proto3.util.equals(RemoveBlockRequest, a, b);
  }
}

/**
 * @generated from message project.RemoveBlockResponse
 */
export class RemoveBlockResponse extends Message<RemoveBlockResponse> {
  /**
   * @generated from field: project.Block block = 1;
   */
  block?: Block;

  constructor(data?: PartialMessage<RemoveBlockResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.RemoveBlockResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "block", kind: "message", T: Block },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RemoveBlockResponse {
    return new RemoveBlockResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RemoveBlockResponse {
    return new RemoveBlockResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RemoveBlockResponse {
    return new RemoveBlockResponse().fromJsonString(jsonString, options);
  }

  static equals(a: RemoveBlockResponse | PlainMessage<RemoveBlockResponse> | undefined, b: RemoveBlockResponse | PlainMessage<RemoveBlockResponse> | undefined): boolean {
    return proto3.util.equals(RemoveBlockResponse, a, b);
  }
}

/**
 * @generated from message project.GetBlocksRequest
 */
export class GetBlocksRequest extends Message<GetBlocksRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  constructor(data?: PartialMessage<GetBlocksRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.GetBlocksRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetBlocksRequest {
    return new GetBlocksRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetBlocksRequest {
    return new GetBlocksRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetBlocksRequest {
    return new GetBlocksRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetBlocksRequest | PlainMessage<GetBlocksRequest> | undefined, b: GetBlocksRequest | PlainMessage<GetBlocksRequest> | undefined): boolean {
    return proto3.util.equals(GetBlocksRequest, a, b);
  }
}

/**
 * @generated from message project.GetBlocksResponse
 */
export class GetBlocksResponse extends Message<GetBlocksResponse> {
  /**
   * @generated from field: repeated project.Block blocks = 1;
   */
  blocks: Block[] = [];

  constructor(data?: PartialMessage<GetBlocksResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.GetBlocksResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "blocks", kind: "message", T: Block, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetBlocksResponse {
    return new GetBlocksResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetBlocksResponse {
    return new GetBlocksResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetBlocksResponse {
    return new GetBlocksResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetBlocksResponse | PlainMessage<GetBlocksResponse> | undefined, b: GetBlocksResponse | PlainMessage<GetBlocksResponse> | undefined): boolean {
    return proto3.util.equals(GetBlocksResponse, a, b);
  }
}

/**
 * @generated from message project.UpdateBlockRequest
 */
export class UpdateBlockRequest extends Message<UpdateBlockRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  /**
   * @generated from field: project.Block block = 2;
   */
  block?: Block;

  constructor(data?: PartialMessage<UpdateBlockRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.UpdateBlockRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "block", kind: "message", T: Block },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UpdateBlockRequest {
    return new UpdateBlockRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UpdateBlockRequest {
    return new UpdateBlockRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UpdateBlockRequest {
    return new UpdateBlockRequest().fromJsonString(jsonString, options);
  }

  static equals(a: UpdateBlockRequest | PlainMessage<UpdateBlockRequest> | undefined, b: UpdateBlockRequest | PlainMessage<UpdateBlockRequest> | undefined): boolean {
    return proto3.util.equals(UpdateBlockRequest, a, b);
  }
}

/**
 * @generated from message project.UpdateBlockResponse
 */
export class UpdateBlockResponse extends Message<UpdateBlockResponse> {
  /**
   * @generated from field: project.Block block = 1;
   */
  block?: Block;

  constructor(data?: PartialMessage<UpdateBlockResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.UpdateBlockResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "block", kind: "message", T: Block },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UpdateBlockResponse {
    return new UpdateBlockResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UpdateBlockResponse {
    return new UpdateBlockResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UpdateBlockResponse {
    return new UpdateBlockResponse().fromJsonString(jsonString, options);
  }

  static equals(a: UpdateBlockResponse | PlainMessage<UpdateBlockResponse> | undefined, b: UpdateBlockResponse | PlainMessage<UpdateBlockResponse> | undefined): boolean {
    return proto3.util.equals(UpdateBlockResponse, a, b);
  }
}

/**
 * @generated from message project.AddEdgeRequest
 */
export class AddEdgeRequest extends Message<AddEdgeRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  /**
   * @generated from field: project.Edge edge = 2;
   */
  edge?: Edge;

  constructor(data?: PartialMessage<AddEdgeRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.AddEdgeRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "edge", kind: "message", T: Edge },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AddEdgeRequest {
    return new AddEdgeRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AddEdgeRequest {
    return new AddEdgeRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AddEdgeRequest {
    return new AddEdgeRequest().fromJsonString(jsonString, options);
  }

  static equals(a: AddEdgeRequest | PlainMessage<AddEdgeRequest> | undefined, b: AddEdgeRequest | PlainMessage<AddEdgeRequest> | undefined): boolean {
    return proto3.util.equals(AddEdgeRequest, a, b);
  }
}

/**
 * @generated from message project.AddEdgeResponse
 */
export class AddEdgeResponse extends Message<AddEdgeResponse> {
  /**
   * @generated from field: project.Edge edge = 1;
   */
  edge?: Edge;

  constructor(data?: PartialMessage<AddEdgeResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.AddEdgeResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "edge", kind: "message", T: Edge },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AddEdgeResponse {
    return new AddEdgeResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AddEdgeResponse {
    return new AddEdgeResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AddEdgeResponse {
    return new AddEdgeResponse().fromJsonString(jsonString, options);
  }

  static equals(a: AddEdgeResponse | PlainMessage<AddEdgeResponse> | undefined, b: AddEdgeResponse | PlainMessage<AddEdgeResponse> | undefined): boolean {
    return proto3.util.equals(AddEdgeResponse, a, b);
  }
}

/**
 * @generated from message project.RemoveEdgeRequest
 */
export class RemoveEdgeRequest extends Message<RemoveEdgeRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  /**
   * @generated from field: string edge_id = 2;
   */
  edgeId = "";

  constructor(data?: PartialMessage<RemoveEdgeRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.RemoveEdgeRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "edge_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RemoveEdgeRequest {
    return new RemoveEdgeRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RemoveEdgeRequest {
    return new RemoveEdgeRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RemoveEdgeRequest {
    return new RemoveEdgeRequest().fromJsonString(jsonString, options);
  }

  static equals(a: RemoveEdgeRequest | PlainMessage<RemoveEdgeRequest> | undefined, b: RemoveEdgeRequest | PlainMessage<RemoveEdgeRequest> | undefined): boolean {
    return proto3.util.equals(RemoveEdgeRequest, a, b);
  }
}

/**
 * @generated from message project.RemoveEdgeResponse
 */
export class RemoveEdgeResponse extends Message<RemoveEdgeResponse> {
  /**
   * @generated from field: project.Edge edge = 1;
   */
  edge?: Edge;

  constructor(data?: PartialMessage<RemoveEdgeResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.RemoveEdgeResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "edge", kind: "message", T: Edge },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RemoveEdgeResponse {
    return new RemoveEdgeResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RemoveEdgeResponse {
    return new RemoveEdgeResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RemoveEdgeResponse {
    return new RemoveEdgeResponse().fromJsonString(jsonString, options);
  }

  static equals(a: RemoveEdgeResponse | PlainMessage<RemoveEdgeResponse> | undefined, b: RemoveEdgeResponse | PlainMessage<RemoveEdgeResponse> | undefined): boolean {
    return proto3.util.equals(RemoveEdgeResponse, a, b);
  }
}

/**
 * @generated from message project.GetEdgesRequest
 */
export class GetEdgesRequest extends Message<GetEdgesRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  constructor(data?: PartialMessage<GetEdgesRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.GetEdgesRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetEdgesRequest {
    return new GetEdgesRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetEdgesRequest {
    return new GetEdgesRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetEdgesRequest {
    return new GetEdgesRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetEdgesRequest | PlainMessage<GetEdgesRequest> | undefined, b: GetEdgesRequest | PlainMessage<GetEdgesRequest> | undefined): boolean {
    return proto3.util.equals(GetEdgesRequest, a, b);
  }
}

/**
 * @generated from message project.GetEdgesResponse
 */
export class GetEdgesResponse extends Message<GetEdgesResponse> {
  /**
   * @generated from field: repeated project.Edge edges = 1;
   */
  edges: Edge[] = [];

  constructor(data?: PartialMessage<GetEdgesResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "project.GetEdgesResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "edges", kind: "message", T: Edge, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetEdgesResponse {
    return new GetEdgesResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetEdgesResponse {
    return new GetEdgesResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetEdgesResponse {
    return new GetEdgesResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetEdgesResponse | PlainMessage<GetEdgesResponse> | undefined, b: GetEdgesResponse | PlainMessage<GetEdgesResponse> | undefined): boolean {
    return proto3.util.equals(GetEdgesResponse, a, b);
  }
}

