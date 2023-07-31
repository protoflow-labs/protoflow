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
import {Node as ProtoNode} from "@/rpc/graph_pb";
import React from "react";
import { useEffect, useState } from "react";
import {NodeEditor} from "@/components/NodeEditor";
import {useUnselect} from "@/components/EditorActions";

const useOverrides = makeStyles({
  card: {},
});

type ActionBarProps = {
  node: ProtoNode | null;
};

export const ActionBar: React.FC<ActionBarProps> = ({node}) => {
  const { selectedNodes } = useSelectedNodes();
  const overrides = useOverrides();

  return (
    <div className="absolute bottom-8 z-10 left-1/2 -translate-x-1/2 w-96">
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
                {node ? <NodeEditor node={node} /> : null}
              </Card>
            </m.div>
          )}
        </AnimatePresence>
      </LazyMotion>
    </div>
  );
}
