import { readFileSync } from "fs";
import functions from './functions/index.js';
import { nodejsService } from "./gen/{{.Runtime}}_connect.js";
import { withReflection } from './reflection/index.js';

export default (router) => {
  return withReflection(
    readFileSync('./gen/image.bin'),
    router.service(nodejsService, functions),
  );
};
