import {Type} from "@/components/jsonViewer/type";

export const JsonViewer = (props: { data: any; }) => {
  return (
    <div className={"json-viewer"}>
      <Type value={props.data} />
    </div>
  )
}
