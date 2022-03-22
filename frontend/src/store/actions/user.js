import * as userReducer from "../reducers/user";
import * as userApi from "../../api/user"
import {setBasic} from "../../api/request";
import * as appReducer from "../reducers/application";



export const signIn = ({username, password}) => dispatch => {
    return userApi.signIn({username, password})
        .then(res => {
            console.log(res)
            if (res?.id) {
                const token = btoa(username+":"+password)
                setBasic(token)
                dispatch(userReducer.login({...res, token}))
            }
            return res
        })
}

export const logout = () => dispatch => {
    dispatch(userReducer.logout())
}

export const checkUser = ({token}) => dispatch => {
    setBasic(token)
    dispatch(userReducer.setToken({token}))
}

export const getPrava = (payload) => dispatch => {
    return userApi.getPrava(payload)
        .then(res => {
            dispatch(userReducer.setPrava(res))
            return res
        })
}

export const setCurrentPodr = (payload) => dispatch => {
    dispatch(userReducer.setCurrentPodr(payload))
}