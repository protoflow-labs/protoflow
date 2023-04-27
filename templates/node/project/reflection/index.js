"use strict";
var __asyncValues = (this && this.__asyncValues) || function (o) {
    if (!Symbol.asyncIterator) throw new TypeError("Symbol.asyncIterator is not defined.");
    var m = o[Symbol.asyncIterator], i;
    return m ? m.call(o) : (o = typeof __values === "function" ? __values(o) : o[Symbol.iterator](), i = {}, verb("next"), verb("throw"), verb("return"), i[Symbol.asyncIterator] = function () { return this; }, i);
    function verb(n) { i[n] = o[n] && function (v) { return new Promise(function (resolve, reject) { v = o[n](v), settle(resolve, reject, v.done, v.value); }); }; }
    function settle(resolve, reject, d, v) { Promise.resolve(v).then(function(v) { resolve({ value: v, done: d }); }, reject); }
};
var __await = (this && this.__await) || function (v) { return this instanceof __await ? (this.v = v, this) : new __await(v); }
var __asyncGenerator = (this && this.__asyncGenerator) || function (thisArg, _arguments, generator) {
    if (!Symbol.asyncIterator) throw new TypeError("Symbol.asyncIterator is not defined.");
    var g = generator.apply(thisArg, _arguments || []), i, q = [];
    return i = {}, verb("next"), verb("throw"), verb("return"), i[Symbol.asyncIterator] = function () { return this; }, i;
    function verb(n) { if (g[n]) i[n] = function (v) { return new Promise(function (a, b) { q.push([n, v, a, b]) > 1 || resume(n, v); }); }; }
    function resume(n, v) { try { step(g[n](v)); } catch (e) { settle(q[0][3], e); } }
    function step(r) { r.value instanceof __await ? Promise.resolve(r.value.v).then(fulfill, reject) : settle(q[0][2], r); }
    function fulfill(value) { resume("next", value); }
    function reject(value) { resume("throw", value); }
    function settle(f, v) { if (f(v), q.shift(), q.length) resume(q[0][0], q[0][1]); }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.withReflection = void 0;
const connect_1 = require("@bufbuild/connect");
const reflection_connect_1 = require("./gen/reflection_connect.js");
const reflection_pb_1 = require("./gen/reflection_pb.js");
const protobuf_1 = require("@bufbuild/protobuf");
const withReflection = (fileDescriptorData, router) => {
    const ds = (0, protobuf_1.createDescriptorSet)(fileDescriptorData);
    router.service(reflection_connect_1.ServerReflection, {
        serverReflectionInfo: reflectionHandler(ds),
    });
    // Some tools haven't been updated to use the v1 release of server reflection so we need
    // to support alpha as well (luckily the schemas are the same)
    router.service(Object.assign(Object.assign({}, reflection_connect_1.ServerReflection), { typeName: "grpc.reflection.v1alpha.ServerReflection" }), {
        serverReflectionInfo: reflectionHandler(ds),
    });
};
exports.withReflection = withReflection;
const reflectionHandler = (ds) => function (reqs) {
    return __asyncGenerator(this, arguments, function* () {
        var _a, e_1, _b, _c;
        try {
            for (var _d = true, reqs_1 = __asyncValues(reqs), reqs_1_1; reqs_1_1 = yield __await(reqs_1.next()), _a = reqs_1_1.done, !_a;) {
                _c = reqs_1_1.value;
                _d = false;
                try {
                    const req = _c;
                    const response = new reflection_pb_1.ServerReflectionResponse({
                        validHost: req.host,
                        originalRequest: req
                    });
                    switch (req.messageRequest.case) {
                        case 'fileByFilename': {
                            response.messageResponse = findFileByFilename(ds, req.messageRequest.value);
                            break;
                        }
                        case 'fileContainingSymbol': {
                            response.messageResponse = findFileContainingSymbol(ds, req.messageRequest.value);
                            break;
                        }
                        case 'listServices': {
                            response.messageResponse = listServices(ds);
                            break;
                        }
                        case 'allExtensionNumbersOfType': {
                            response.messageResponse = errorResponse('Not currently implemented', connect_1.Code.Unimplemented);
                            break;
                        }
                        case 'fileContainingExtension': {
                            response.messageResponse = errorResponse('Not currently implemented', connect_1.Code.Unimplemented);
                            break;
                        }
                        default: {
                            response.messageResponse = notFoundErrorResponse(req.messageRequest.value || '');
                        }
                    }
                    yield yield __await(response);
                }
                finally {
                    _d = true;
                }
            }
        }
        catch (e_1_1) { e_1 = { error: e_1_1 }; }
        finally {
            try {
                if (!_d && !_a && (_b = reqs_1.return)) yield __await(_b.call(reqs_1));
            }
            finally { if (e_1) throw e_1.error; }
        }
    });
};
const listServices = (ds) => {
    return {
        case: 'listServicesResponse',
        value: new reflection_pb_1.ListServiceResponse({
            service: [...ds.services.keys()].map(serviceName => new reflection_pb_1.ServiceResponse({
                name: serviceName
            }))
        }),
    };
};
const findFileByFilename = (ds, filename) => {
    const sanitizedFilename = filename.replace(/\.proto$/, '');
    const file = ds.files.find(f => f.name === sanitizedFilename);
    if (file) {
        return fileDescriptorResponse(file);
    }
    return notFoundErrorResponse(filename);
};
// Rough port of https://github.com/protocolbuffers/protobuf-go/blob/808c66411fe76c839a68168654228446fbdc1ecf/reflect/protoregistry/registry.go#L222
const findFileContainingSymbol = (ds, name) => {
    let prefix = name;
    let suffix = '';
    while (prefix != '') {
        const lookup = lookupDescriptor(ds, prefix);
        if (lookup) {
            const { type, descriptor } = lookup;
            if (descriptor.typeName === name) {
                return fileDescriptorResponse(descriptor.file);
            }
            switch (type) {
                case 'message': {
                    // Handle potential sub message match
                    break;
                }
                case 'service': {
                    const method = descriptor.methods.find(m => m.name === suffix);
                    if (method) {
                        return fileDescriptorResponse(method.parent.file);
                    }
                }
            }
        }
        const parts = prefix.split('.');
        const trailing = parts.pop() || '';
        suffix = suffix ? (suffix + '.' + trailing) : trailing;
        prefix = parts.join('.');
    }
    return notFoundErrorResponse(name);
};
const lookupDescriptor = (ds, name) => {
    const service = ds.services.get(name);
    if (service) {
        return {
            type: 'service',
            descriptor: service,
        };
    }
    const descEnum = ds.enums.get(name);
    if (descEnum) {
        return {
            type: 'enum',
            descriptor: descEnum,
        };
    }
    const message = ds.messages.get(name);
    if (message) {
        return {
            type: 'message',
            descriptor: message,
        };
    }
    const extension = ds.extensions.get(name);
    if (extension) {
        return {
            type: 'extension',
            descriptor: extension,
        };
    }
};
const fileDescriptorResponse = (file) => ({
    case: 'fileDescriptorResponse',
    value: new reflection_pb_1.FileDescriptorResponse({
        fileDescriptorProto: [file.proto.toBinary()],
    })
});
const errorResponse = (errorMessage, errorCode) => ({
    case: 'errorResponse',
    value: new reflection_pb_1.ErrorResponse({
        errorMessage,
        errorCode,
    }),
});
const notFoundErrorResponse = (name) => errorResponse(`Could not find: ${name}`, connect_1.Code.NotFound);
