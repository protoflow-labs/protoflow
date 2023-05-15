import {Type} from "@/components/jsonViewer/type";

export const JsonViewer = (props: { data: any; }) => {
  return (
    <div>
      <Type value={props.data} />
    </div>
  )
}
