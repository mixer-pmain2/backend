import * as appReducer from "../reducers/application"
import * as appApi from "../../api/spr"

export const enableLoading = () => dispatch =>
    dispatch(appReducer.loadingEnable())

export const disableLoading = () => dispatch =>
    dispatch(appReducer.loadingDisable())

export const getSprPodr = () => dispatch => {
    return appApi.getSprPodr()
        .then(r => dispatch(appReducer.setSprPodr(r)))
}

export const getSprPrava = () => dispatch => {
    return appApi.getSprPrava()
        .then(r => dispatch(appReducer.setSprPrava(r)))
}