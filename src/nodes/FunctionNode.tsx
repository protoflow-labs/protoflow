import { useCallback } from "react";
import { Handle, NodeProps, Position } from "reactflow";
import { HiCircleStack, HiCodeBracket } from "react-icons/hi2";

const handleStyle = {
  background: "blue",
};
type FunctionNodeProps = NodeProps<FunctionData>;

type FunctionData = {
  name: string;
  language: string;
};

export function FunctionNode(props: FunctionNodeProps) {
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
            <label htmlFor="text">Language:</label>
            <select
              id="language"
              name="language"
              onChange={(e) => {
                props.data.language = e.target.value;
              }}
            >
              <option>Node.js</option>
              <option>Python</option>
              <option>Go</option>
            </select>
          </>
        ) : (
          <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
            <HiCodeBracket />
            {data.name} ({data.language})
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
