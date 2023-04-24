import { HiCheckCircle } from "react-icons/hi2";
import { Connection, Handle, NodeProps, Position } from "reactflow";

const handleStyle = {
  background: "green",
};

type ValidatorNodeProps = NodeProps<ValidatorData>;

type ValidatorData = {
  name: string;
};

const handleValidConnection = (connection: Connection) => {
  return true;
};

export function ValidatorNode(props: ValidatorNodeProps) {
  const { data, selected } = props;

  return (
    <>
      <Handle
        type="source"
        position={Position.Bottom}
        isValidConnection={handleValidConnection}
      />
      <div className="flex flex-col p-4 rounded-tl-md rounded-br-md bg-white text-black text-xs">
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
          </>
        ) : (
          <div className="flex items-center gap-8">
            <HiCheckCircle />
            {data.name}
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
      <Handle
        type="source"
        position={Position.Right}
        id="onError"
        style={{ background: "red" }}
      />
    </>
  );
}
