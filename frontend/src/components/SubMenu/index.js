import React from "react";


const SubMenu = ({tabs, onChange, curTab, style}) => {
    style = {...style, cursor: "pointer"}
    return <ul className="nav nav-tabs" style={style}>
        {
            tabs.map((v, i) => <li className="nav-item" key={i}>
                <a
                    className={`nav-link ${curTab === v.id && "active"}`}
                    aria-current="page"
                    onClick={_ => onChange(v.id)}
                >
                    {v.title}
                </a>
            </li>)
        }

    </ul>
}

export default SubMenu