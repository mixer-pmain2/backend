import React from "react";


const SubMenu = ({tabs, onChange, curTab, style}) => {
  style = {...style, cursor: "pointer"}
  return <ul className="nav nav-tabs" style={{...style}}>
    {
      tabs.map((v, i) => {
        const isActive = curTab === v.id
        return <li className="nav-item" key={i}>
          <a
            className={`nav-link ${isActive && "active"}`}
            aria-current="page"
            onClick={_ => onChange(v.id)}
          >
            {/*{v.img &&*/}
            {/*  <img style={{*/}
            {/*    width: 9,*/}
            {/*    marginRight: 5,*/}
            {/*    filter: "invert(48%) sepia(79%) saturate(2476%) hue-rotate(86deg) brightness(118%) contrast(119%)"*/}
            {/*  }}*/}
            {/*       src={v.img.default} alt={v.title} title={v.title}*/}
            {/*  />}*/}
            {v.title}
          </a>
        </li>
      })
    }

  </ul>
}

export default SubMenu