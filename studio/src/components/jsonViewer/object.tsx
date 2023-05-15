import React, { useState } from 'react'
import {Type} from "@/components/jsonViewer/type";

const KeyExpander = (props: { keyName: any; value: any; startExpanded?: boolean }) => {
  const { keyName, value, startExpanded = false } = props
  const [expanded, setExpanded] = useState(startExpanded)

  if (typeof value !== 'object' && !Array.isArray(value)) {
    return (
      <>
        <span>{keyName}: </span>
        <Type value={value} />
      </>
    )
  }

  return (
    <>
      <span
        className="rsjv-expander"
        onClick={() => setExpanded((prev: boolean) => !prev)}
      >
        [{expanded ? '-' : '+'}]{' '}
      </span>
      <span>{keyName}: </span>
      {expanded && <Type value={value} />}
    </>
  )
}

export const ObjectType = (props: { value: any }) => {
  const { value } = props

  return (
    <ul>
      {Object.keys(value || {}).map(k => (
        <li key={k}>
          <KeyExpander keyName={k} value={value[k]} />
        </li>
      ))}
    </ul>
  )
}
