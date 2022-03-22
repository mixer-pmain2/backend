import React, {useState} from "react";
import {connect} from "react-redux";

import * as appActions from "../../store/actions/application"

import Notify, {notifyType} from "../../components/Notify";

import {SideLoading} from "../../components/Loading";
import Layout from "../Layout";

const TestPage = ({dispatch}) => {
    const [isLoading, setIsLoading] = useState(true)
    return <Layout>
        <div>
            <button onClick={_ => Notify(notifyType.INFO, "info")()}>info</button>
            <button onClick={_ => Notify(notifyType.SUCCESS, "success")()}>success</button>
            <button onClick={_ => Notify(notifyType.ERROR, "error")()}>error</button>
            <button onClick={_ => Notify(notifyType.WARNING, "warning")()}>warning</button>
        </div>
        <div>
            <button onClick={_ => setIsLoading(true)}>on</button>
            <button onClick={_ => setIsLoading(false)}>off</button>
            <SideLoading isLoading={isLoading}/>
        </div>
        <div>
            <button onClick={_ => dispatch(appActions.enableLoading())}>on</button>
            <button onClick={_ => dispatch(appActions.disableLoading())}>off</button>
            <SideLoading isLoading={isLoading}/>
        </div>
    </Layout>
}

export default connect(state => ({}))(TestPage)