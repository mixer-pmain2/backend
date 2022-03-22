import React from "react";
import {connect} from "react-redux";

import * as userAction from "../../store/actions/user"
import Layout from "../Layout";

const SignInPage = ({dispatch}) => {

    const handleSubmit = (e) => {
        e.preventDefault()
        const [username, password] = e.target
        dispatch(userAction.signIn({
            username: username.value,
            password: password.value
        }))
    }

    return <div style={{marginTop: 20}}>
        <h4>Авторизация</h4>
        <form onSubmit={handleSubmit}>
            <div className="input-group mb-3">
                <span className="input-group-text">Логин</span>
                <input type="text" className="form-control" placeholder="Username" aria-label="Username"
                       aria-describedby="username" name="username" id="username"/>
            </div>
            <div className="input-group mb-3">
                <span className="input-group-text">Пароль</span>
                <input type="password" className="form-control" placeholder="Password" aria-label="Password"
                       aria-describedby="password" name="password" id="password"/>
            </div>
            <input type="submit" className="form-control btn btn-outline-primary" value="Вход"/>
        </form>
    </div>
}

export default connect((state) => ({}))(SignInPage)
