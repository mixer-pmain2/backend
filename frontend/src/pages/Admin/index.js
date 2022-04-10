import React, {useState} from "react";
import {connect} from "react-redux";

import Layout from "../Layout";
import SubMenu from "components/SubMenu";
import Test from "./Test";
import UserParams from "./UserParams";
import DiagnoseTree from "../Patient/components/DiagnoseTree";

const tabs = [
    {
        id: 0,
        title: "Пользователь"
    },
    {
        id: 1,
        title: "Diags"
    },
    {
        id: 99,
        title: "Test"
    },
]

const AdminPage = (props) => {
    const [curTab, setCurTab] = useState(tabs[0].id)

    return <Layout>
        <SubMenu tabs={tabs} curTab={curTab} onChange={setCurTab} style={{marginBottom: 35}}/>
        {curTab === tabs[0].id && <UserParams/>}
        {curTab === tabs[1].id && <DiagnoseTree/>}
        {curTab === tabs[2].id && <Test {...props}/>}
    </Layout>
}

export default connect(state => ({
    application: state.application,
    user: state.user
}))(AdminPage)