import React from "react";
import {connect} from "react-redux";
import {useNavigate} from "react-router-dom";
import {linkDict} from "../routes";

const AuthRequire = ({children, podr, dispatch, user}) => {
    const navigate = useNavigate()

    console.log("middleware Auth")

    if (!user.isAuth) {
        // navigate(linkDict.signin)
    }

    return <>{children}</>
}

const AuthRequireState = connect(state => ({
    user: state.user
}))(AuthRequire)

const middlewareAuthRequire = (next, params) => {
    return <AuthRequireState {...params}>
        {next}
    </AuthRequireState>
}

export default middlewareAuthRequire