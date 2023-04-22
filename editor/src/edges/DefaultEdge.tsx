import { EdgeProps, getBezierPath } from "reactflow";

type DefaultEdgeData = {
  async?: boolean;
};

export default function DefaultEdge({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  style = {},
  data = { async: false },
  markerEnd,
}: EdgeProps<DefaultEdgeData>) {
  const [edgePath] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  return (
    <>
      <path
        id={id}
        style={{ ...style, stroke: data.async ? "green" : "gray" }}
        className="react-flow__edge-path"
        d={edgePath}
        markerEnd={markerEnd}
        onDoubleClick={() => {
          data.async = !data.async;
        }}
      />
    </>
  );
}
