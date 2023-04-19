import { useCallback } from "react";
import { Handle, NodeProps, Position } from "reactflow";
import { HiCircleStack } from "react-icons/hi2";

const handleStyle = {
  background: "#555",
};
type EntityNodeProps = NodeProps<EntityData>;

type EntityData = {
  name: string;
  table: string;
};

export function EntityNode(props: EntityNodeProps) {
  const { data, selected } = props;
  return (
    <>
      <Handle type="source" position={Position.Bottom} />
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          padding: "16px",
          borderRadius: "6px",
          background: "white",
          color: "black",
        }}
      >
        {selected ? (
          <>
            <label htmlFor="text">Name:</label>
            <input
              id="text"
              name="text"
              onChange={(e) => {
                console.log(e.target.value);
                props.data.name = e.target.value;
              }}
              className="nodrag"
              value={data.name}
            />
            <label htmlFor="text">Table:</label>
            <input
              id="text"
              name="text"
              onChange={(e) => {
                props.data.table = e.target.value;
              }}
              className="nodrag"
              value={data.table}
            />
          </>
        ) : (
          <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
            <HiCircleStack />
            {data.name} ({data.table})
          </div>
        )}
      </div>
      <Handle type="target" position={Position.Top} id="a" />
      <Handle
        type="source"
        position={Position.Bottom}
        id="b"
        style={handleStyle}
      />
    </>
  );
}
