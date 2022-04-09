import * as userReducer from "../reducers/user";
import * as userApi from "../../api/user"
import {setBasic} from "../../api/request";



export const signIn = ({username, password}) => dispatch => {
  const token = btoa(unescape(encodeURIComponent(username + ":" + password)))
  // const token = btoa(username + ":" + password)
  console.log(token)
  return userApi.signIn({token})
    .then(res => {

      if (res?.id) {
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
  return dispatch(userReducer.setToken({token}))
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

export const getUch = (payload) => dispatch => {
  return userApi.getUch(payload)
    .then(res => {
      dispatch(userReducer.setUch(res))
      return res
    })
}