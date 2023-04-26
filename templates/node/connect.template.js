import { NodeService } from "./gen/node_connect.js";
import methods from './methods';

export default (router) => {
  return router.service(NodeService, methods);
};
