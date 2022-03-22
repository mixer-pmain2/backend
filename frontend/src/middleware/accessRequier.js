import React from "react";
import {connect} from "react-redux";

import NotAccessedPage from "../pages/NotAccessed";

const AccessRequire = (props) => {
    const {children, dispatch, user, params} = props

    console.log("middleware Access", user.podr)
    console.log(user)

    console.log(props)

    const curPodr = user?.podr
    const curPrava = user?.prava[curPodr]

    const checkAccess = (accessList) => {
        for (let i=0; i < accessList.length; i++) {
            const accessPodr = accessList[i].podr
            const accessPrava = accessList[i].prava
            console.log(accessPodr, accessPrava, curPodr, curPrava)
            if (curPodr === accessPodr && (curPrava & accessPrava) > 0)
                return true
        }
        return false
    }

    if (!checkAccess(params.access)) {
        return <NotAccessedPage/>
    }


    return <>{children}</>
}

const AccessRequireState = connect(state => ({
    user: state.user
}))(AccessRequire)

const middlewareAccessRequire = (next, params) => {
    return <AccessRequireState params={params}>
        {next}
    </AccessRequireState>
}

export default middlewareAccessRequire