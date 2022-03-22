import React from "react"
import SubMenu from "../../../../components/SubMenu";

export const AmbTabs = ["Посещения", "Учет", ]

const Menu = ({onChange, curTab}) => {

    return <SubMenu
        tabs={AmbTabs}
        onChange={onChange}
        curTab={curTab}
    />
}

export default Menu