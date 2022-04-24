import {API, paramsToUrlQuery, request} from "./request";


export const findByFio = ({fio}) => {
    console.log(fio)
    const [lname, fname, sname] = fio.split(" ")
    const url = API + `/patient/find/?lname=${lname||""}&fname=${fname||""}&sname=${sname||""}`
    return request("GET", url, {}, {})
}

export const findByID = ({id}) => {
    const url = API + `/patient/${id}/`
    return request("GET", url, {}, {})
}

export const findUchet = ({id}) => {
    const url = API + `/patient/${id}/uchet/`
    return request("GET", url, {}, {})
}

export const getVisits = (payload) => {
    const url = API + `/patient/${payload.id}/visit/?`+paramsToUrlQuery(payload)
    return request("GET", url, {}, {})
}

export const getHospital = (payload) => {
    const url = API + `/patient/${payload.id}/hospital/?`+paramsToUrlQuery(payload)
    return request("GET", url, {}, {})
}

export const newVisit = (payload) => {
    const url = API + `/patient/${payload.patientId}/visit/`
    return request("POST", url, {}, {...payload})
}

export const newProf = (payload) => {
    const url = API + `/patient/prof/`
    return request("POST", url, {}, {...payload})
}
