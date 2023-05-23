import {
  Button,
  Dialog, DialogActions,
  DialogBody,
  DialogContent,
  DialogSurface,
  DialogTitle,
  DialogTrigger,
  FluentProvider,
  webDarkTheme
} from "@fluentui/react-components";
import {ReactFlowProvider} from "reactflow";
import {HotkeysProvider} from "react-hotkeys-hook/src/HotkeysProvider";
import ProjectProvider from "@/providers/ProjectProvider";
import {EditorProvider} from "@/providers/EditorProvider";
import {Toolbar} from "@/components/Toolbar";
import {Toaster} from "react-hot-toast";
import {AppRoutes} from "@/routes";
import {BrowserRouter} from "react-router-dom";
import {ErrorBoundary} from "react-error-boundary";
import {ErrorInfo, useState} from "react";

const Fallback: React.FC<{ error: Error, resetErrorBoundary: () => void }> = ({error, resetErrorBoundary}) => {
  const [open, setOpen] = useState(true);
  return (
    <Dialog open={open}>
      <DialogSurface>
        <DialogBody>
          <DialogTitle>Unhandled Error</DialogTitle>
          <DialogContent>
            <h4>{error.message.toString()}</h4>
            <pre>{error.stack}</pre>
          </DialogContent>
          <DialogActions>
            <DialogTrigger disableButtonEnhancement>
              <Button appearance="secondary" onClick={() => resetErrorBoundary()}>Close</Button>
            </DialogTrigger>
          </DialogActions>
        </DialogBody>
      </DialogSurface>
    </Dialog>
  );
}

export default function App() {
  return (
      <FluentProvider theme={webDarkTheme}>
        <ErrorBoundary
          FallbackComponent={Fallback}
        >
        <ReactFlowProvider>
          <HotkeysProvider initiallyActiveScopes={["editor"]}>
            <ProjectProvider>
              <EditorProvider>
                  <Toolbar/>
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
