import React, {useState} from "react";


const TreeItem = ({title, value, onSelect}) => {
  const [toggleHover, setToggleHover] = useState(false)
  const style = toggleHover ? {backgroundColor: "whitesmoke"} : {}

  return <div
    style={style}
    onMouseEnter={() => setToggleHover(true)}
    onMouseLeave={() => setToggleHover(false)}
    onClick={_ => onSelect(value)}
  >
    <button className="btn" style={{width: 40}}/>
    <span>{value} {title}</span>
  </div>
}

const TreeNode = ({node, onOpen, onSelect}) => {
  const [data, setData] = useState([])
  const [isOpen, setIsOpen] = useState(false)

  const handleClick = () => {
    setIsOpen(!isOpen)
    if (!isOpen && data.length === 0) onOpen(node.value, setData)
  }


  return <div style={{cursor: "pointer"}}>
    <a className="btn" style={{width: 40}} onClick={handleClick}>{isOpen ? "-" : "+"}</a>
    <span onClick={handleClick}>{node.value} {node.title}</span>
    {isOpen && data.length > 0 && <div>
      <TreeView data={data} getData={(v, setData) => onOpen(v, setData)} onSelect={onSelect}/>
    </div>}
  </div>
}

export const TreeView = ({data, getData, onSelect}) => {

  return <div>
    {
      data.map((v, i) => <div key={i} style={{marginLeft: 25}}>
          {
            v.isNode ?
              <TreeNode node={v} onOpen={getData} onSelect={onSelect}/> :
              <TreeItem value={v.value} title={v.title} onSelect={onSelect}/>
          }
        </div>
      )
    }
  </div>
}