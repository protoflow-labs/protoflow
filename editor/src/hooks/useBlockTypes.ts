import { BucketBlock } from "@/components/blocks/BucketBlock";
import { FunctionBlock } from "@/components/blocks/FunctionBlock";
import { InputBlock } from "@/components/blocks/InputBlock";
import { QueryBlock } from "@/components/blocks/QueryBlock";
import { QueueBlock } from "@/components/blocks/QueueBlock";
import { useMemo } from "react";
import {CollectionBlock} from "@/components/blocks/CollectionBlock";

export function useBlockTypes() {
  const nodeTypes = useMemo(
    () => ({
      "protoflow.collection": CollectionBlock,
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
