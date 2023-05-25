import { connectNodeAdapter } from "@bufbuild/connect-node";
import * as http2 from "http2";
import routes from "./connect.js";

const PORT = 8086;
const server = connectNodeAdapter({ routes });

http2.createServer(server).listen(PORT, () => {
  console.log(`Listening on port ${PORT}`);
});
