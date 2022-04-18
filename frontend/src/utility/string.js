export const formatDate = (s = "", format="dd.mm.YYYY") => {
    // s = 1998-11-18
    const [year, month, day] = s.split("-")
    let res = []
    const _format = format.split(".")
    _format.map((v) => {
        if (v === "dd") res.push(day)
        if (v === "mm") res.push(month)
        if (v === "YYYY") res.push(year)
    })
    return res.join(".")
}

export const formatDateToInput = (d) => {
    let month = (d.getMonth()+1)
    if (month < 10) month = '0'+month
    let day = d.getDate()
    if (day < 10) day = '0'+day
    let date = d.getFullYear() + '-' + month + '-' + day
    return date
}

export function capitalizeFirstLetter(string) {
    return string?.charAt(0).toUpperCase() + string?.slice(1);
}

export const shorty = (text, len) =>
    (text || "").substring(0, len)+(text?.length>len?"...":"")