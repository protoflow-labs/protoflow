import { BucketBlock } from "@/blocks/BucketBlock";
import { EntityBlock } from "@/blocks/EntityBlock";
import { FunctionBlock } from "@/blocks/FunctionBlock";
import { InputBlock } from "@/blocks/InputBlock";
import { QueryBlock } from "@/blocks/QueryBlock";
import { QueueBlock } from "@/blocks/QueueBlock";
import { useMemo } from "react";

export function useNodeTypes() {
  const nodeTypes = useMemo(
    () => ({
      "protoflow.entity": EntityBlock,
      "protoflow.function": FunctionBlock,
      "protoflow.input": InputBlock,
      "protoflow.query": QueryBlock,
      "protoflow.queue": QueueBlock,
      "protoflow.bucket": BucketBlock,
    }),
    []
  );

  return { nodeTypes };
}
