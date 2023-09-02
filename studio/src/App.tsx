import {
  FluentProvider,
  webDarkTheme
} from "@fluentui/react-components";
import {ReactFlowProvider} from "reactflow";
import {HotkeysProvider} from "react-hotkeys-hook/src/HotkeysProvider";
import ProjectProvider from "@/providers/ProjectProvider";
import {EditorProvider} from "@/providers/EditorProvider";
import {Toaster} from "react-hot-toast";
import {BrowserRouter, useRoutes} from "react-router-dom";
import {ErrorBoundary} from "react-error-boundary";
import {FallbackError} from "@/components/FallbackError";
import "react-chat-elements/dist/main.css"
import { initializeIcons } from '@fluentui/react/lib/Icons';
import Home from "@/pages";

initializeIcons();

const AppRoutes = () => {
    const commonRoutes = [{
        path: '/studio',
        element: <Home />
    }];

    const element = useRoutes([...commonRoutes]);

    return <>{element}</>;
};

export default function App() {
  return (
      <FluentProvider theme={webDarkTheme}>
        <ErrorBoundary
          FallbackComponent={FallbackError}
        >
          <ReactFlowProvider>
          <HotkeysProvider initiallyActiveScopes={["editor"]}>
            <ProjectProvider>
              <EditorProvider>
                  <BrowserRouter>
                    <AppRoutes/>
                  </BrowserRouter>
                  <Toaster/>
              </EditorProvider>
            </ProjectProvider>
          </HotkeysProvider>
        </ReactFlowProvider>
        </ErrorBoundary>
      </FluentProvider>
  )
}
