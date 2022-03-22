import React from "react";


const SubMenu = ({tabs, onChange, curTab}) => {
    return <ul className="nav nav-tabs">
        {
            tabs.map((v, i) => <li className="nav-item" key={i} onClick={_ => onChange(v)}>
                <a className={`nav-link ${curTab === v && "active"}`}
                   aria-current="page"
                   href="#"
                >{v}</a>
            </li>)
        }

    </ul>
}

export default SubMenu