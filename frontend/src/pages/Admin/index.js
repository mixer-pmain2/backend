import React, {useState} from "react";
import {connect} from "react-redux";

import Layout from "../Layout";
import SubMenu from "../../components/SubMenu";
import Test from "./Test";
import UserParams from "./UserParams";

const tabs = [
    {
        id: 0,
        title: "Test"
    },
    {
        id: 1,
        title: "User params"
    }
]

const AdminPage = (props) => {
    const [curTab, setCurTab] = useState(tabs[0].id)

    return <Layout>
        <SubMenu tabs={tabs} curTab={curTab} onChange={setCurTab} style={{marginBottom: 35}}/>
        {curTab === tabs[0].id && <Test {...props}/>}
        {curTab === tabs[1].id && <UserParams/>}
    </Layout>
}

export default connect(state => ({
    application: state.application,
    user: state.user
}))(AdminPage)