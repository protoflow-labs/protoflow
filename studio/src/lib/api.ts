import {
  createConnectTransport,
} from "@bufbuild/connect-web";
import { GenerateService } from "@/rpc/generate_connect";
import { ProjectService } from "@/rpc/project_connect";
import {createPromiseClient} from "@bufbuild/connect";

const transport = createConnectTransport({
  baseUrl: "/api",
});

export const projectService = createPromiseClient(ProjectService, transport);
export const generateService = createPromiseClient(GenerateService, transport);