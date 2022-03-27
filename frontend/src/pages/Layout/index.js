import React, {useEffect} from "react"
import {connect} from "react-redux";

import SignInPage from "../Signin";

import {checkUser, logout} from "../../store/actions/user";

import NotificationContainer from "react-notifications/lib/NotificationContainer";
import {SideLoading} from "../../components/Loading";

import NavMenu from "./NavMenu";


function Layout({children, user, application, dispatch, patient}) {
    useEffect(() => {
        let action = false
        if (!action) {
            if (user?.token) {
                dispatch(checkUser({token: user?.token}))
            }
        }


        return () => {
            action = true
        }
    }, [dispatch, user?.token])

    const handleLogout = () => dispatch(logout())

    return <div className={`container`}>
        {
            user?.isAuth ? <div>
                <div className="mb-5">
                    <NavMenu onLogout={handleLogout} user={user} app={application} patient={patient}/>
                </div>
                {children}
            </div> : <SignInPage/>
        }
        <NotificationContainer/>
        <SideLoading isLoading={application.loading}/>
    </div>
}

export default connect((state) => ({
    user: state.user,
    application: state.application,
    patient: state.patient
}))(Layout)