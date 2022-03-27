import axios from "axios";

import * as userActions from "../store/actions/user"
import Notify, {notifyType} from "../components/Notify";

export const API = "http://localhost:80/api/v0"

export let basicAuth = "";

export const setBasic = (token) => {
    basicAuth = "Basic "+token
}


export const request = (method, url, headers = {}, body = {}) => {
    console.log("request", method, url)
    method = method.toUpperCase()
    return axios({
        method,
        url,
        headers: {
            'Authorization': basicAuth,
            'Content-Type': "text/plain",
            ...headers
        },
        data: body
    })
        .then(res => {
            console.log(res)
            if (res.status === 200) {
                return res.data
            }
            return {}
        })
        .catch(reqError)
}

export const reqError = (err) => {
    console.log(err.response.status)

    switch (err.response.status) {
        case 401:
            Notify(notifyType.ERROR, "Ошибка авторизации")()
            break
        default:
            Notify(notifyType.ERROR, err.message)()
    }

}

