import React from "react";
import {connect} from "react-redux";

import NotAccessedPage from "../pages/NotAccessed";

import {isAccessed} from "../configs/access";

const AccessRequire = (props) => {
    const {children, dispatch, user, params} = props

    console.log("middleware Access")

    const userUnit = user?.unit
    const userAccess = user?.access?.[userUnit]

    if (!isAccessed(params.access, userUnit, userAccess)) {
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