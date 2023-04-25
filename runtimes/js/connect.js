import { Code, ConnectError } from "@bufbuild/connect";
import { RunService } from "./gen/run_connect.js";

export default (registry) => (router) => {
  return router.service(RunService, {
    // implements rpc Run
    async run(req) {
      const func = registry[req.functionName];
      if (!func) {
        throw new ConnectError("Function not found", Code.InvalidArgument);
      }

      const output = JSON.stringify(await func(parse(req.input)));

      return {
        output,
      };
    },
  });
};

function parse(input) {
  try {
    return JSON.parse(input);
  } catch {
    return input;
  }
}
