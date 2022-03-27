import React from "react";
import {connect} from "react-redux";

import Layout from "../pages/Layout";
import SignInPage from "../pages/Signin";

const AuthRequire = (props) => {
    const {children, dispatch, user} = props

    console.log("middleware Auth")

    if (!user.isAuth) {
        return <Layout>
            <SignInPage/>
        </Layout>
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