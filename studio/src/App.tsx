import {FluentProvider, webDarkTheme} from "@fluentui/react-components";
import {ReactFlowProvider} from "reactflow";
import {HotkeysProvider} from "react-hotkeys-hook/src/HotkeysProvider";
import ProjectProvider from "@/providers/ProjectProvider";
import {EditorProvider} from "@/providers/EditorProvider";
import {Toolbar} from "@/components/Toolbar";
import {Toaster} from "react-hot-toast";
import {AppRoutes} from "@/routes";
import {BrowserRouter} from "react-router-dom";

export default function App() {
  return (
    <FluentProvider theme={webDarkTheme}>
      <ReactFlowProvider>
        <HotkeysProvider initiallyActiveScopes={["editor"]}>
          <ProjectProvider>
            <EditorProvider>
              <Toolbar />
              <BrowserRouter>
                <AppRoutes />
              </BrowserRouter>
              <Toaster />
            </EditorProvider>
          </ProjectProvider>
        </HotkeysProvider>
      </ReactFlowProvider>
    </FluentProvider>
  )
}
