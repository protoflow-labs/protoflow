import {
  createConnectTransport,
  createPromiseClient,
} from "@bufbuild/connect-web";
import { ProjectService } from "../../rpc/project_connect";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

export const projectService = createPromiseClient(ProjectService, transport);
