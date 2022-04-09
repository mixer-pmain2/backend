import {API, request} from "./request";


export const signIn = ({token}) => {
    const url = API + "/auth/signin/"
    const headers = {
        "Authorization": "Basic "+token
    }
    return request("GET", url, headers, {})
}

export const getPrava = ({id}) => {
    const url = API + `/user/${id}/prava/`
    return request('GET', url, {}, {})
}

export const getUser = ({id}) => {
    const url = API + `/user/${id}/`
    return request('GET', url, {}, {})
}

export const getUch = ({id}) => {
    const url = API + `/user/${id}/uch/`
    return request('GET', url, {}, {})
}