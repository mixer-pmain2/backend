import React from "react"
import SubMenu from "../../../../components/SubMenu";

const Menu = ({onChange, curTab, tabs, style}) => {

    return <SubMenu
        tabs={tabs}
        onChange={onChange}
        curTab={curTab}
        style={style}
    />
}

export default Menu