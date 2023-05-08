import {
  Dialog,
  DialogTrigger,
  DialogSurface,
  DialogTitle,
  DialogBody,
  DialogActions,
  DialogContent,
  Button, Field, Input,
} from "@fluentui/react-components";
import * as React from "react";
import {FC} from "react";
import {useForm} from "react-hook-form";
import {useUnselect} from "@/components/EditorActions";
import {projectService} from "@/lib/api";
import {useProjectContext} from "@/providers/ProjectProvider";
import {toast} from "react-hot-toast";

interface AddResourceDialogProps {
  open: boolean;
  close: () => void;
}

export const AddResourceDialog: FC<AddResourceDialogProps> = ({open, close}) => {
  const { project, loadResources, loadingResources } = useProjectContext();
  const onCancel = useUnselect();

  const { register, handleSubmit, watch } = useForm({
    values: {
      url: '',
      name: '',
    },
  });

  const onSubmit = async (data: any) => {
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
      resource: {
        name: data.name,
        type: {
          case: 'grpcService',
          value: {
            host: data.url,
          }
        }
      }
    })
    toast.success('Resource added');
    await loadResources();
    close();
  };

  const values = watch();

  return (
    <Dialog open={open} onOpenChange={close}>
      <DialogSurface>
        <DialogBody>
          <DialogTitle>Add Resource</DialogTitle>
          <DialogContent>
            Addd an existing GRPC resource to the project. The GRPC resource must have the reflection API enabled for importing to work correctly.
            <form>
              <div className="flex flex-col gap-2 p-3">
                <Field label="Name" required>
                  <Input value={values.name} {...register('name')} />
                </Field>
                <Field label="URL" required>
                  <Input value={values.url} {...register('url')} />
                </Field>
              </div>
            </form>
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
