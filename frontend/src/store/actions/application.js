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