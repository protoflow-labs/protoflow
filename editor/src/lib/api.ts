import {
  createConnectTransport,
  createPromiseClient,
} from "@bufbuild/connect-web";
import { GenerateService } from "../../rpc/generate_connect";
import { ProjectService } from "../../rpc/project_connect";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

export const projectService = createPromiseClient(ProjectService, transport);
export const generateService = createPromiseClient(GenerateService, transport);
