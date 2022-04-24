import React, {useCallback, useEffect} from "react"
import {connect} from "react-redux";

import SignInPage from "../Signin";

import {checkUser, logout} from "../../store/actions/user";

import NotificationContainer from "react-notifications/lib/NotificationContainer";
import Loading, {SideLoading, typeLoading} from "../../components/Loading";

import NavMenu from "./NavMenu";
import Progress from "../../components/Progress";


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

    const navMenu = useCallback((user, application, patient) =>
        <NavMenu onLogout={handleLogout} user={user} app={application} patient={patient}/>,
        []
    )
    return <div className={`container`}>
        {
            user?.isAuth ? <div>
                <div className="mb-5">
                    {navMenu(user, application, patient)}
                </div>
                {children}
            </div> : <SignInPage/>
        }
        <NotificationContainer/>
        <SideLoading isLoading={application.loading}/>
        <div style={{position: "fixed", bottom: 35, right: 35}}>
            <Loading type={typeLoading.INFO} isLoading={application?.loadingList?.length > 0}/>
        </div>
        {application?.loadingList?.length > 0 && <Progress style={{position: "fixed", top: 0, left: 0, width: "100%"}}/>}
        <footer style={{height: 150}}>

        </footer>
    </div>
}

export default connect((state) => ({
    user: state.user,
    application: state.application,
    patient: state.patient
}))(Layout)