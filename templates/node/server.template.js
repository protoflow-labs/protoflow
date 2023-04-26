import { connectNodeAdapter } from "@bufbuild/connect-node";
import * as http from "http";
import routes from "./connect.js";

const PORT = 8086;
const server = connectNodeAdapter({ routes });

http.createServer(server).listen(PORT, () => {
  console.log(`Listening on port ${PORT}`);
});
