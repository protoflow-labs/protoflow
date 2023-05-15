import React from 'react'
import {Type} from "@/components/jsonViewer/type";

export const ArrayType = (props: { value: any; }) => {
  const { value } = props

  return (
    <ul>
      {value.map((v: any, i: any) => (
        <li key={i}>
          <Type value={v} />
        </li>
      ))}
    </ul>
  )
}
