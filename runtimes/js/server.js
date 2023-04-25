import { connectNodeAdapter } from "@bufbuild/connect-node";
import * as http from "http";
import routes from "./connect.js";
import { register } from "./register.js";

register().then((registry) => {
  const PORT = 8086;
  const server = connectNodeAdapter({ routes: routes(registry) });

  http.createServer(server).listen(PORT, () => {
    console.log(`Listening on port ${PORT}`);
  });
});
