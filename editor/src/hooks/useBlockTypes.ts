import { BucketBlock } from "@/components/blocks/BucketBlock";
import { EntityBlock } from "@/components/blocks/EntityBlock";
import { FunctionBlock } from "@/components/blocks/FunctionBlock";
import { InputBlock } from "@/components/blocks/InputBlock";
import { QueryBlock } from "@/components/blocks/QueryBlock";
import { QueueBlock } from "@/components/blocks/QueueBlock";
import { useMemo } from "react";

export function useBlockTypes() {
  const blockTypes = useMemo(
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

  return { blockTypes };
}
