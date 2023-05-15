import React, { useState } from 'react'

export const ValueType = (props: { value: any }) => {
  const { value } = props
  const [showFull, setShowFull] = useState(false)

  // shorten long strings until clicked
  if (typeof value === 'string' && !showFull) {
    return (
      <span className="rsjv-value" onClick={() => setShowFull(true)}>
        {value.length > 100 ? `${value.substring(0, 100)}...` : value}
      </span>
    )
  }

  return (
    <span className="rsjv-value" onClick={() => setShowFull(false)}>
      {value}
    </span>
  )
}
