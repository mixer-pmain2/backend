import {API, request} from "./request";


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