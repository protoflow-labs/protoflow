import { BucketBlock } from "@/components/blocks/BucketBlock";
import { CollectionBlock } from "@/components/blocks/CollectionBlock";
import { FunctionBlock } from "@/components/blocks/FunctionBlock";
import { InputBlock } from "@/components/blocks/InputBlock";
import { QueryBlock } from "@/components/blocks/QueryBlock";
import { QueueBlock } from "@/components/blocks/QueueBlock";
import { RESTBlock } from "@/components/blocks/RESTBlock";
import { useMemo } from "react";

export function useBlockTypes() {
  const nodeTypes = useMemo(
    () => ({
      "protoflow.collection": CollectionBlock,
      "protoflow.function": FunctionBlock,
      "protoflow.input": InputBlock,
      "protoflow.query": QueryBlock,
      "protoflow.queue": QueueBlock,
      "protoflow.bucket": BucketBlock,
      "protoflow.rest": RESTBlock,
    }),
    []
  );

  return { nodeTypes };
}
