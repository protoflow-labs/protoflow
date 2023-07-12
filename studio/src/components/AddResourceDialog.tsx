import {
  Dialog,
  DialogTrigger,
  DialogSurface,
  DialogTitle,
  DialogBody,
  DialogActions,
  DialogContent,
  Button, Field, Input, Textarea, Select,
} from "@fluentui/react-components";
import * as React from "react";
import {FC, useState} from "react";
import {useForm} from "react-hook-form";
import {useUnselect} from "@/components/EditorActions";
import {projectService} from "@/lib/api";
import {useProjectContext} from "@/providers/ProjectProvider";
import {toast} from "react-hot-toast";
import {GRPCService, Resource } from "@/rpc/resource_pb";
import {AnyMessage, MessageType} from "@bufbuild/protobuf";

interface AddResourceDialogProps {
  open: boolean;
  close: () => void;
}

let configLookup: Record<string, MessageType> = {};
Resource.fields.list().forEach((field) => {
  if (field.oneof && field.oneof.name === 'type') {
    if (field.kind === 'message') {
      configLookup[field.name] = field.T;
    }
  }
})

function snakeToCamel(s: string): string {
  const words = s.split('_');
  const camelWords = words.map((word, index) => {
    if (index === 0) {
      return word.toLowerCase();
    } else {
      return capitalize(word);
    }
  });
  return camelWords.join('');
}

function capitalize(word: string): string {
  return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
}

export const AddResourceDialog: FC<AddResourceDialogProps> = ({open, close}) => {
  const { project, loadResources, loadingResources } = useProjectContext();
  const onCancel = useUnselect();

  const [configTypeName, setConfigTypeName] = useState<string|null>(null);
  const configType = configTypeName ? configLookup[configTypeName] : null;

  const {watch, setValue, register, handleSubmit, control} = useForm({
    values: {
      name: "",
      config: {},
    },
  });
  const values = watch();

  const onSubmit = async (data: any) => {
    if (!configType) {
      toast.error('No config type selected');
      return;
    }
    if (!project) {
      toast.error('No project loaded');
      return;
    }
    if (loadingResources) {
      console.log('resources are already being loaded');
      return;
    }
    // TODO breadchris support different resource types
    await projectService.createResource({
      projectId: project.id,
      resource: new Resource({
        name: data.name,
        type: {
          //@ts-ignore
          case: snakeToCamel(configTypeName),
          value: configType.fromJson(data.config),
        }
      })
    })
    toast.success('Resource added');
    await loadResources();
    onCancel();
  };

  const resourceChanged = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setValue('name', '')
    setValue('config', {});
    setConfigTypeName(e.target.value);
  }

  const fields = configType ? configType.fields.list() : null;

  return (
    <Dialog open={open} onOpenChange={close}>
      <DialogSurface>
        <DialogBody>
          <DialogTitle>Add Resource</DialogTitle>
          <DialogContent>
            <Select onChange={resourceChanged}>
              {Object.keys(configLookup).map((key) => {
                return <option key={key} value={key}>{key}</option>
              })}
            </Select>
            <Field label="Name" required>
              <Input value={values.name} {...register('name')} />
            </Field>
            {fields && fields.map((field) => {
              return (
                  <Field label={field.name} key={field.name}>
                    {/* @ts-ignore */}
                    <Textarea value={values.config[field.name] || ''} {...register(`config.${field.name}`)} />
                  </Field>
              )
            })}
          </DialogContent>
          <DialogActions>
            <DialogTrigger disableButtonEnhancement>
              <Button appearance="secondary" disabled={loadingResources}>Close</Button>
            </DialogTrigger>
            <Button appearance="primary" onClick={handleSubmit(onSubmit)} disabled={loadingResources}>Add</Button>
          </DialogActions>
        </DialogBody>
      </DialogSurface>
    </Dialog>
  );
};
