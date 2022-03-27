import * as patientApi from "../../api/patient"
import * as patientReducer from "../../store/reducers/patient"

export const findByFio = ({fio}) => dispatch =>
    patientApi.findByFio({fio})
        .then(res => {
            return res
        })

export const findById = ({id}) => dispatch =>
    patientApi.findByID({id})
        .then(res => {
            return res
        })

export const select = (item) => dispatch =>
    dispatch(patientReducer.select(item))

export const reset = () => dispatch =>
    dispatch(patientReducer.reset())

export const getUchet = (payload) => dispatch =>
    patientApi.findUchet(payload)
        .then(res => {
            dispatch(patientReducer.setUchet(res))
            return res
        })