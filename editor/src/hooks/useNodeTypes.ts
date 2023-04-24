import { BucketBlock } from "@/blocks/BucketBlock";
import { EntityBlock } from "@/blocks/EntityBLock";
import { FunctionBlock } from "@/blocks/FunctionBlock";
import { InputBlock } from "@/blocks/InputBlock";
import { QueryBlock } from "@/blocks/QueryBlock";
import { QueueBlock } from "@/blocks/QueueBlock";
import { useMemo } from "react";

export function useNodeTypes() {
  const nodeTypes = useMemo(
    () => ({
      entity: EntityBlock,
      function: FunctionBlock,
      message: InputBlock,
      query: QueryBlock,
      queue: QueueBlock,
      bucket: BucketBlock,
    }),
    []
  );

  return { nodeTypes };
}
