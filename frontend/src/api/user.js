import {API, request} from "./request";



export const signIn = ({username, password}) => {
    const url = API + "/auth/signin/"
    const headers = {
        "Authorization": "Basic "+btoa(username+":"+password)
    }
    return request("GET", url, headers, {})
}

export const getPrava = ({id}) => {
    const url = API + `/user/${id}/prava/`
    return request('GET', url, {}, {})
}