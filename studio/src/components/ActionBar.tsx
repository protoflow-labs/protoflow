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

  if (selectedNodes.length === 0) {
    return null;
  }

  return (
    <div>
      {node ? <NodeEditor node={node} /> : null}
    </div>
  );
}
