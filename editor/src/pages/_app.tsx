import { Toolbar } from "@/components/Toolbar";
import { EditorProvider } from "@/providers/EditorProvider";
import ProjectProvider from "@/providers/ProjectProvider";
import "@/styles/globals.css";
import { FluentProvider, webDarkTheme } from "@fluentui/react-components";
import type { AppProps } from "next/app";
import { Toaster } from "react-hot-toast";
import { HotkeysProvider } from "react-hotkeys-hook";
import { ReactFlowProvider } from "reactflow";

import "reactflow/dist/style.css";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <FluentProvider theme={webDarkTheme}>
      <ReactFlowProvider>
        <HotkeysProvider initiallyActiveScopes={["editor"]}>
          <ProjectProvider>
            <EditorProvider>
              <Toolbar />
              <Component {...pageProps} />
              <Toaster />
            </EditorProvider>
          </ProjectProvider>
        </HotkeysProvider>
      </ReactFlowProvider>
    </FluentProvider>
  );
}
