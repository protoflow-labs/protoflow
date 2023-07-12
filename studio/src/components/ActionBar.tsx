import { useSelectedNodes } from "@/hooks/useSelectedNodes";
import { useProjectContext } from "@/providers/ProjectProvider";
import {
  Button,
  Card,
  CardHeader,
  Text,
  makeStyles,
} from "@fluentui/react-components";
import { AnimatePresence, LazyMotion, domAnimation, m } from "framer-motion";
import React from "react";
import { useEffect, useState } from "react";
import {JsonViewer} from "@/components/jsonViewer";

const useOverrides = makeStyles({
  card: {},
});

export function ActionBar() {
  const { project, runWorkflow } = useProjectContext();
  const { selectedNodes } = useSelectedNodes();
  const [loading] = useState(false);
  const overrides = useOverrides();

  const onRun = async () => {
    if (!project) return;

    await runWorkflow(selectedNodes[0]);
  };

  return (
    <div className="absolute bottom-8 z-10 left-1/2 -translate-x-1/2">
      <LazyMotion features={domAnimation}>
        <AnimatePresence>
          {selectedNodes.length === 1 && (
            <m.div
              initial={{ opacity: 0, translateY: 10, scale: 0.95 }}
              animate={{ opacity: 1, translateY: 0, scale: 1 }}
              exit={{ opacity: 0, translateY: 10, scale: 0.95 }}
              transition={{ duration: 0.1 }}
              className="flex flex-col gap-4"
            >
              <Card className={overrides.card}>
                <Button onClick={onRun} disabled={loading}>
                  Run Workflow
                </Button>
              </Card>
            </m.div>
          )}
        </AnimatePresence>
      </LazyMotion>
    </div>
  );
}
