import React, {useEffect} from "react"
import {connect} from "react-redux";
import {Link} from "react-router-dom";

import SignInPage from "../Signin";

import {checkUser, logout} from "../../store/actions/user";

import {linkDict} from "../../routes";
import NotificationContainer from "react-notifications/lib/NotificationContainer";
import {SideLoading} from "../../components/Loading";


const NavMenu = ({onLogout, user, app}) => {
    return <nav className="navbar navbar-expand-lg navbar-light bg-light">
        <div className="container-fluid">
            <Link to={linkDict.start} className="navbar-brand">КПБ</Link>
            <div className="collapse navbar-collapse" id="navbarSupportedContent">
                <ul className="navbar-nav me-auto mb-2 mb-lg-0">
                    <li className="nav-item">
                        <Link to={linkDict.findPatient} className="nav-link">{app?.spr?.podr[user?.podr]}</Link>
                    </li>
                </ul>
            </div>
            <div>
                <span className="m-2 navbar-text">Пользователь: {user?.id}</span>
                <button className="btn btn-primary" onClick={onLogout}>Выход</button>
            </div>
        </div>
    </nav>
}

function Layout({children, user, application, dispatch}) {
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
                    <NavMenu onLogout={handleLogout} user={user} app={application}/>
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
    application: state.application
}))(Layout)