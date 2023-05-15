import React from 'react'
import {ObjectType} from "@/components/jsonViewer/object";
import {ValueType} from "@/components/jsonViewer/value";
import {ArrayType} from "@/components/jsonViewer/array";

export const Type = (props: { [x: string]: any; value: any }) => {
  const { value, ...rest } = props

  if (Array.isArray(value)) {
    return <ArrayType value={value} {...rest} />
  }

  if (typeof value === 'object') {
    return <ObjectType value={value} {...rest} />
  }

  return <ValueType value={value} {...rest} />
}
