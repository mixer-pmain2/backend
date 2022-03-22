import axios from "axios";

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
    console.log(err)
    Notify(notifyType.ERROR, err.getMessage)()
}

