import { FunctionData } from "@/blocks/FunctionBlock";
import { Divider, Input, Label, Select } from "@fluentui/react-components";
import { Node } from "reactflow";
import { EditorActions } from "../EditorActions";

export function FunctionEditor(props: { node: Node<FunctionData> }) {
  return (
    <div className="flex flex-col gap-2 p-3">
      <div className="flex flex-col">
        <Label htmlFor="entityName">Name</Label>
        <Input
          id="entityName"
          defaultValue={props.node.data.name || ""}
          onChange={(e) => {
            props.node.data.name = e.target.value;
          }}
        />
      </div>
      <div className="flex flex-col">
        <Label htmlFor="entityLanguage">Language</Label>
        <Select
          onChange={(e) => {
            // props.node.data.language = e.currentTarget.value;
          }}
        >
          <option>Go</option>
          <option>Node.js</option>
          <option>Python</option>
        </Select>
      </div>
      <Divider />
      <EditorActions />
    </div>
  );
}
