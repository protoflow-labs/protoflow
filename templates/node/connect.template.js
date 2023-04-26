import { NodeService } from "./gen/node_connect.js";
import methods from './methods.js';

export default (router) => {
  return router.service(NodeService, methods);
};
