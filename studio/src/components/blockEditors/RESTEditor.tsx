import {
  Button,
  Divider,
  Dropdown,
  Field,
  Input,
  Option,
} from "@fluentui/react-components";
import { useForm } from "react-hook-form";
import { HiOutlineTrash, HiPlus } from "react-icons/hi2";
import { Node } from "reactflow";
import { REST } from "@/rpc/block_pb";
import { EditorActions, useUnselect } from "../EditorActions";
import { RESTData } from "@/components/blocks/RESTBlock";

export function RESTEditor({ node }: { node: Node<RESTData> }) {
  const onCancel = useUnselect();
  const { watch, setValue, register, handleSubmit } = useForm({
    values: {
      name: node.data.name || "",
      config: {
        method: node.data.config.rest?.method || "GET",
        path: node.data.config.rest?.path || "",
      } as REST,
      headers: Object.entries(node.data.config.rest?.headers || {}).map(
        ([key, value]) => ({ key, value })
      ),
    },
  });

  const onSubmit = (data: any) => {
    node.data.name = data.name;

    if (!node.data.config.rest) {
      node.data.config.rest = {
        method: "",
        path: "",
        headers: {},
      };
    }

    node.data.config.rest.method = data.config.method;
    node.data.config.rest.path = data.config.path;
    node.data.config.rest.headers = (
      data.headers as { key: string; value: string }[]
    ).reduce((obj, acc) => {
      obj[acc.key] = acc.value;
      return obj;
    }, {} as Record<string, string>);

    onCancel();
  };

  const values = watch();

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="flex flex-col gap-2 p-3">
        <Field label="Name" required>
          <Input value={values.name} {...register("name")} />
        </Field>
        <Field label="Method" required>
          <Dropdown
            value={"GET"}
            onOptionSelect={(_, data) => {
              setValue("config.method", data.optionValue || "");
            }}
          >
            <Option>GET</Option>
            <Option>PATCH</Option>
            <Option>POST</Option>
            <Option>PUT</Option>
            <Option>DELETE</Option>
          </Dropdown>
        </Field>
        <Field label="Path" required>
          <Input value={values.config.path} {...register("config.path")} />
        </Field>
        <div className="flex flex-col">
          <Field label="Headers">
            {values.headers.map((header, index) => {
              return (
                <div key={index} className="flex items-center gap-2 mb-2">
                  <Input
                    value={header.key}
                    {...register(`headers.${index}.key`)}
                  />
                  <Input
                    value={header.value}
                    {...register(`headers.${index}.value`)}
                  />

                  <Button
                    icon={<HiOutlineTrash className="h-4 w-4" />}
                    onClick={() => {
                      setValue(
                        "headers",
                        values.headers.filter((_, i) => i !== index)
                      );
                    }}
                  />
                </div>
              );
            })}
          </Field>

          <Button
            icon={<HiPlus className="w-4 h-4" />}
            onClick={() => {
              setValue(
                "headers",
                values.headers.concat({ key: "", value: "" })
              );
            }}
          >
            Add Header
          </Button>
        </div>
        <Divider />
        <EditorActions />
      </div>
    </form>
  );
}
