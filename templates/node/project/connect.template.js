import { readFileSync } from "fs";
import functions from './functions/index.js';
import { nodeService } from "./gen/node.service_connect.js";
import { withReflection } from './reflection/index.js';

export default (router) => {
  return withReflection(
    readFileSync('./gen/image.bin'),
    router.service(nodeService, functions),
  );
};
