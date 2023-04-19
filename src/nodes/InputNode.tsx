import {
  ChangeEvent,
  KeyboardEvent,
  useEffect,
  useState
} from "react";
import { useForm } from "react-hook-form";
import { HiPencilSquare, HiPlay, HiXMark } from "react-icons/hi2";
import { Handle, NodeProps, Position } from "reactflow";

const handleStyle = {
  background: "#555",
};

type InputNodeProps = NodeProps<InputData>;

type Field = {
  name: string;
  type: "string" | "number" | "boolean";
  required?: boolean;
};

type InputData = {
  name: string;
  fields: Field[];
  lastUpdate: number;
};

type InputProps = {
  defaultValue?: string;
  onKeyUp?: (e: KeyboardEvent<HTMLInputElement>) => void;
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
};

function Input(props: InputProps) {
  return (
    <input
      id="text"
      name="text"
      type="text"
      autoComplete="off"
      onKeyUp={props.onKeyUp}
      onChange={props.onChange}
      className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
      defaultValue={props.defaultValue}
    />
  );
}

export function InputNode(props: InputNodeProps) {
  const [editing, setEditing] = useState(false);
  const [addingField, setAddingField] = useState(false);
  const { data, selected } = props;
  const { watch, setValue } = useForm({
    defaultValues: props.data,
  });

  const values = watch();

  useEffect(() => {
    if (!selected) {
      setEditing(false);
    }
  }, [selected]);

  useEffect(() => {
    data.fields = values.fields;
    data.name = values.name;
  }, [values]);

  return (
    <>
      <Handle type="source" position={Position.Bottom} />
      <div
        onClick={() => {
          setEditing(true);
        }}
      >
        {selected && editing ? (
          <>
            <label htmlFor="text">Name:</label>
            <Input
              defaultValue={data.name}
              onChange={(e) => {
                setValue("name", e.target.value);
              }}
            />
            {values.fields?.map((field) => (
              <p>
                {field.name}: {field.type}
                <button
                  onClick={() => {
                    setValue(
                      "fields",
                      values.fields.filter((f) => f.name !== field.name)
                    );
                  }}
                >
                  <HiXMark />
                </button>
              </p>
            ))}
            {addingField ? (
              <Input
                onKeyUp={(e) => {
                  if (e.key === "Enter") {
                    if (!values.fields) {
                      setValue("fields", []);
                    }

                    setValue(
                      "fields",
                      values.fields.concat({
                        name: e.currentTarget.value,
                        type: "string",
                      })
                    );

                    setAddingField(false);
                  }
                  if (e.key === "Escape") {
                    setAddingField(false);
                  }
                }}
              />
            ) : (
              <button
                onClick={() => {
                  setAddingField(true);
                }}
              >
                Add field
              </button>
            )}
          </>
        ) : (
          <div className="flex items-center gap-2">
            <HiPencilSquare />
            {data.name}

            <button
              style={{
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                background: "green",
                borderRadius: "50%",
                width: "40px",
                height: "40px",
                position: "absolute",
                top: -20,
                right: -20,
              }}
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
              }}
            >
              <HiPlay />
            </button>
          </div>
        )}
      </div>
      <Handle
        type="source"
        position={Position.Bottom}
        id="b"
        style={handleStyle}
      />
    </>
  );
}
