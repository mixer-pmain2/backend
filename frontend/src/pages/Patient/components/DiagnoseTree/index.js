import React, {useEffect, useState} from "react";

import * as apiSpr from "api/spr"


const TreeItem = ({title, value, onSelect}) => {
  return <div>
    <button className="btn" style={{width: 40}} onClick={_ => onSelect(value)}/>
    <span onClick={_ => onSelect(value)}>{value} {title}</span>
  </div>
}

const TreeNode = ({node, onOpen}) => {
  const [data, setData] = useState([])
  const [isOpen, setIsOpen] = useState(false)

  const handleClick = () => {
    setIsOpen(!isOpen)
    if (!isOpen && data.length === 0) onOpen(node.value, setData)
  }

  return <div style={{cursor: "pointer"}}>
    <button className="btn" style={{width: 40}} onClick={handleClick}>{isOpen ? "-" : "+"}</button>
    <span onClick={handleClick}>{node.value} {node.title}</span>
    {isOpen && data.length > 0 && <div>
      <TreeView data={data} getData={(v, setData) => onOpen(v, setData)}/>
    </div>}
  </div>
}

const TreeView = ({data, getData}) => {

  return <div>
    {
      data.map((v, i) => <div key={i} style={{marginLeft: 25}}>
          {
            v.isNode ?
              <TreeNode node={v} onOpen={getData}/> :
              <TreeItem value={v.value} title={v.title} onSelect={console.log}/>
          }
        </div>
      )
    }
  </div>
}

const DiagnoseTree = () => {
  const [diags, setDiags] = useState([])

  const mapper = (data) => data.map((v, i) => ({
    title: v.title,
    value: v.diag,
    isNode: v.haveChildren
  }))

  const getDiag = (diag = "", mapper, setData) => {
    apiSpr.getSprDiag({diag})
      .then(res => {
        setData(mapper(res))
      })
  }

  useEffect(() => {
    getDiag("", mapper, setDiags)
  }, [])

  console.log(diags)
  return <div>
    <TreeView data={diags} getData={(v, setData) => getDiag(v, mapper, setData)}/>
  </div>
}

export default DiagnoseTree