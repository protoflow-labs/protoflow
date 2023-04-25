import { ConnectRouter } from "@bufbuild/connect";
import { Service } from "./gen/service";

export default (router: ConnectRouter) => {
  router.service(Service, {
    // implements rpc Say
    async execute(req) {
      return {
        sentence: `You said: ${req.sentence}`,
      };
    },
  });
};
