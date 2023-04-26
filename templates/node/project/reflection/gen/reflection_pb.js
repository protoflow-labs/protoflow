"use strict";
// Copyright 2016 The gRPC Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
Object.defineProperty(exports, "__esModule", { value: true });
exports.ErrorResponse = exports.ServiceResponse = exports.ListServiceResponse = exports.ExtensionNumberResponse = exports.FileDescriptorResponse = exports.ServerReflectionResponse = exports.ExtensionRequest = exports.ServerReflectionRequest = void 0;
const protobuf_1 = require("@bufbuild/protobuf");
/**
 * The message sent by the client when calling ServerReflectionInfo method.
 *
 * @generated from message grpc.reflection.v1.ServerReflectionRequest
 */
class ServerReflectionRequest extends protobuf_1.Message {
    constructor(data) {
        super();
        /**
         * @generated from field: string host = 1;
         */
        this.host = "";
        /**
         * To use reflection service, the client should set one of the following
         * fields in message_request. The server distinguishes requests by their
         * defined field and then handles them using corresponding methods.
         *
         * @generated from oneof grpc.reflection.v1.ServerReflectionRequest.message_request
         */
        this.messageRequest = { case: undefined };
        protobuf_1.proto3.util.initPartial(data, this);
    }
    static fromBinary(bytes, options) {
        return new ServerReflectionRequest().fromBinary(bytes, options);
    }
    static fromJson(jsonValue, options) {
        return new ServerReflectionRequest().fromJson(jsonValue, options);
    }
    static fromJsonString(jsonString, options) {
        return new ServerReflectionRequest().fromJsonString(jsonString, options);
    }
    static equals(a, b) {
        return protobuf_1.proto3.util.equals(ServerReflectionRequest, a, b);
    }
}
ServerReflectionRequest.runtime = protobuf_1.proto3;
ServerReflectionRequest.typeName = "grpc.reflection.v1.ServerReflectionRequest";
ServerReflectionRequest.fields = protobuf_1.proto3.util.newFieldList(() => [
    { no: 1, name: "host", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "file_by_filename", kind: "scalar", T: 9 /* ScalarType.STRING */, oneof: "message_request" },
    { no: 4, name: "file_containing_symbol", kind: "scalar", T: 9 /* ScalarType.STRING */, oneof: "message_request" },
    { no: 5, name: "file_containing_extension", kind: "message", T: ExtensionRequest, oneof: "message_request" },
    { no: 6, name: "all_extension_numbers_of_type", kind: "scalar", T: 9 /* ScalarType.STRING */, oneof: "message_request" },
    { no: 7, name: "list_services", kind: "scalar", T: 9 /* ScalarType.STRING */, oneof: "message_request" },
]);
exports.ServerReflectionRequest = ServerReflectionRequest;
/**
 * The type name and extension number sent by the client when requesting
 * file_containing_extension.
 *
 * @generated from message grpc.reflection.v1.ExtensionRequest
 */
class ExtensionRequest extends protobuf_1.Message {
    constructor(data) {
        super();
        /**
         * Fully-qualified type name. The format should be <package>.<type>
         *
         * @generated from field: string containing_type = 1;
         */
        this.containingType = "";
        /**
         * @generated from field: int32 extension_number = 2;
         */
        this.extensionNumber = 0;
        protobuf_1.proto3.util.initPartial(data, this);
    }
    static fromBinary(bytes, options) {
        return new ExtensionRequest().fromBinary(bytes, options);
    }
    static fromJson(jsonValue, options) {
        return new ExtensionRequest().fromJson(jsonValue, options);
    }
    static fromJsonString(jsonString, options) {
        return new ExtensionRequest().fromJsonString(jsonString, options);
    }
    static equals(a, b) {
        return protobuf_1.proto3.util.equals(ExtensionRequest, a, b);
    }
}
ExtensionRequest.runtime = protobuf_1.proto3;
ExtensionRequest.typeName = "grpc.reflection.v1.ExtensionRequest";
ExtensionRequest.fields = protobuf_1.proto3.util.newFieldList(() => [
    { no: 1, name: "containing_type", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "extension_number", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
]);
exports.ExtensionRequest = ExtensionRequest;
/**
 * The message sent by the server to answer ServerReflectionInfo method.
 *
 * @generated from message grpc.reflection.v1.ServerReflectionResponse
 */
class ServerReflectionResponse extends protobuf_1.Message {
    constructor(data) {
        super();
        /**
         * @generated from field: string valid_host = 1;
         */
        this.validHost = "";
        /**
         * The server sets one of the following fields according to the message_request
         * in the request.
         *
         * @generated from oneof grpc.reflection.v1.ServerReflectionResponse.message_response
         */
        this.messageResponse = { case: undefined };
        protobuf_1.proto3.util.initPartial(data, this);
    }
    static fromBinary(bytes, options) {
        return new ServerReflectionResponse().fromBinary(bytes, options);
    }
    static fromJson(jsonValue, options) {
        return new ServerReflectionResponse().fromJson(jsonValue, options);
    }
    static fromJsonString(jsonString, options) {
        return new ServerReflectionResponse().fromJsonString(jsonString, options);
    }
    static equals(a, b) {
        return protobuf_1.proto3.util.equals(ServerReflectionResponse, a, b);
    }
}
ServerReflectionResponse.runtime = protobuf_1.proto3;
ServerReflectionResponse.typeName = "grpc.reflection.v1.ServerReflectionResponse";
ServerReflectionResponse.fields = protobuf_1.proto3.util.newFieldList(() => [
    { no: 1, name: "valid_host", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "original_request", kind: "message", T: ServerReflectionRequest },
    { no: 4, name: "file_descriptor_response", kind: "message", T: FileDescriptorResponse, oneof: "message_response" },
    { no: 5, name: "all_extension_numbers_response", kind: "message", T: ExtensionNumberResponse, oneof: "message_response" },
    { no: 6, name: "list_services_response", kind: "message", T: ListServiceResponse, oneof: "message_response" },
    { no: 7, name: "error_response", kind: "message", T: ErrorResponse, oneof: "message_response" },
]);
exports.ServerReflectionResponse = ServerReflectionResponse;
/**
 * Serialized FileDescriptorProto messages sent by the server answering
 * a file_by_filename, file_containing_symbol, or file_containing_extension
 * request.
 *
 * @generated from message grpc.reflection.v1.FileDescriptorResponse
 */
class FileDescriptorResponse extends protobuf_1.Message {
    constructor(data) {
        super();
        /**
         * Serialized FileDescriptorProto messages. We avoid taking a dependency on
         * descriptor.proto, which uses proto2 only features, by making them opaque
         * bytes instead.
         *
         * @generated from field: repeated bytes file_descriptor_proto = 1;
         */
        this.fileDescriptorProto = [];
        protobuf_1.proto3.util.initPartial(data, this);
    }
    static fromBinary(bytes, options) {
        return new FileDescriptorResponse().fromBinary(bytes, options);
    }
    static fromJson(jsonValue, options) {
        return new FileDescriptorResponse().fromJson(jsonValue, options);
    }
    static fromJsonString(jsonString, options) {
        return new FileDescriptorResponse().fromJsonString(jsonString, options);
    }
    static equals(a, b) {
        return protobuf_1.proto3.util.equals(FileDescriptorResponse, a, b);
    }
}
FileDescriptorResponse.runtime = protobuf_1.proto3;
FileDescriptorResponse.typeName = "grpc.reflection.v1.FileDescriptorResponse";
FileDescriptorResponse.fields = protobuf_1.proto3.util.newFieldList(() => [
    { no: 1, name: "file_descriptor_proto", kind: "scalar", T: 12 /* ScalarType.BYTES */, repeated: true },
]);
exports.FileDescriptorResponse = FileDescriptorResponse;
/**
 * A list of extension numbers sent by the server answering
 * all_extension_numbers_of_type request.
 *
 * @generated from message grpc.reflection.v1.ExtensionNumberResponse
 */
class ExtensionNumberResponse extends protobuf_1.Message {
    constructor(data) {
        super();
        /**
         * Full name of the base type, including the package name. The format
         * is <package>.<type>
         *
         * @generated from field: string base_type_name = 1;
         */
        this.baseTypeName = "";
        /**
         * @generated from field: repeated int32 extension_number = 2;
         */
        this.extensionNumber = [];
        protobuf_1.proto3.util.initPartial(data, this);
    }
    static fromBinary(bytes, options) {
        return new ExtensionNumberResponse().fromBinary(bytes, options);
    }
    static fromJson(jsonValue, options) {
        return new ExtensionNumberResponse().fromJson(jsonValue, options);
    }
    static fromJsonString(jsonString, options) {
        return new ExtensionNumberResponse().fromJsonString(jsonString, options);
    }
    static equals(a, b) {
        return protobuf_1.proto3.util.equals(ExtensionNumberResponse, a, b);
    }
}
ExtensionNumberResponse.runtime = protobuf_1.proto3;
ExtensionNumberResponse.typeName = "grpc.reflection.v1.ExtensionNumberResponse";
ExtensionNumberResponse.fields = protobuf_1.proto3.util.newFieldList(() => [
    { no: 1, name: "base_type_name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "extension_number", kind: "scalar", T: 5 /* ScalarType.INT32 */, repeated: true },
]);
exports.ExtensionNumberResponse = ExtensionNumberResponse;
/**
 * A list of ServiceResponse sent by the server answering list_services request.
 *
 * @generated from message grpc.reflection.v1.ListServiceResponse
 */
class ListServiceResponse extends protobuf_1.Message {
    constructor(data) {
        super();
        /**
         * The information of each service may be expanded in the future, so we use
         * ServiceResponse message to encapsulate it.
         *
         * @generated from field: repeated grpc.reflection.v1.ServiceResponse service = 1;
         */
        this.service = [];
        protobuf_1.proto3.util.initPartial(data, this);
    }
    static fromBinary(bytes, options) {
        return new ListServiceResponse().fromBinary(bytes, options);
    }
    static fromJson(jsonValue, options) {
        return new ListServiceResponse().fromJson(jsonValue, options);
    }
    static fromJsonString(jsonString, options) {
        return new ListServiceResponse().fromJsonString(jsonString, options);
    }
    static equals(a, b) {
        return protobuf_1.proto3.util.equals(ListServiceResponse, a, b);
    }
}
ListServiceResponse.runtime = protobuf_1.proto3;
ListServiceResponse.typeName = "grpc.reflection.v1.ListServiceResponse";
ListServiceResponse.fields = protobuf_1.proto3.util.newFieldList(() => [
    { no: 1, name: "service", kind: "message", T: ServiceResponse, repeated: true },
]);
exports.ListServiceResponse = ListServiceResponse;
/**
 * The information of a single service used by ListServiceResponse to answer
 * list_services request.
 *
 * @generated from message grpc.reflection.v1.ServiceResponse
 */
class ServiceResponse extends protobuf_1.Message {
    constructor(data) {
        super();
        /**
         * Full name of a registered service, including its package name. The format
         * is <package>.<service>
         *
         * @generated from field: string name = 1;
         */
        this.name = "";
        protobuf_1.proto3.util.initPartial(data, this);
    }
    static fromBinary(bytes, options) {
        return new ServiceResponse().fromBinary(bytes, options);
    }
    static fromJson(jsonValue, options) {
        return new ServiceResponse().fromJson(jsonValue, options);
    }
    static fromJsonString(jsonString, options) {
        return new ServiceResponse().fromJsonString(jsonString, options);
    }
    static equals(a, b) {
        return protobuf_1.proto3.util.equals(ServiceResponse, a, b);
    }
}
ServiceResponse.runtime = protobuf_1.proto3;
ServiceResponse.typeName = "grpc.reflection.v1.ServiceResponse";
ServiceResponse.fields = protobuf_1.proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
]);
exports.ServiceResponse = ServiceResponse;
/**
 * The error code and error message sent by the server when an error occurs.
 *
 * @generated from message grpc.reflection.v1.ErrorResponse
 */
class ErrorResponse extends protobuf_1.Message {
    constructor(data) {
        super();
        /**
         * This field uses the error codes defined in grpc::StatusCode.
         *
         * @generated from field: int32 error_code = 1;
         */
        this.errorCode = 0;
        /**
         * @generated from field: string error_message = 2;
         */
        this.errorMessage = "";
        protobuf_1.proto3.util.initPartial(data, this);
    }
    static fromBinary(bytes, options) {
        return new ErrorResponse().fromBinary(bytes, options);
    }
    static fromJson(jsonValue, options) {
        return new ErrorResponse().fromJson(jsonValue, options);
    }
    static fromJsonString(jsonString, options) {
        return new ErrorResponse().fromJsonString(jsonString, options);
    }
    static equals(a, b) {
        return protobuf_1.proto3.util.equals(ErrorResponse, a, b);
    }
}
ErrorResponse.runtime = protobuf_1.proto3;
ErrorResponse.typeName = "grpc.reflection.v1.ErrorResponse";
ErrorResponse.fields = protobuf_1.proto3.util.newFieldList(() => [
    { no: 1, name: "error_code", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 2, name: "error_message", kind: "scalar", T: 9 /* ScalarType.STRING */ },
]);
exports.ErrorResponse = ErrorResponse;
//# sourceMappingURL=reflection_pb.js.map
