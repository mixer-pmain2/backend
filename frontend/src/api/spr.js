import {API, request} from "./request";

export const getSprPodr = () => {
    const url = API + "/spr/podr/"
    return request('GET', url, {}, {})
}

export const getSprPrava = () => {
    const url = API + "/spr/prava/"
    return request('GET', url, {}, {})
}

export const getSprVisit = () => {
    const url = API + "/spr/visit/"
    return request('GET', url, {}, {})
}