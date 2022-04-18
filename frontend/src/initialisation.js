import React from "react";
import {connect} from "react-redux";

import * as appActions from "./store/actions/application"
import {setBasic} from "./api/request";


// Не добавлять состояния иначе будет часто вызываться
const Initialisation = ({dispatch}) => {

    dispatch(appActions.loadingReset())
    dispatch(appActions.getSprPodr())
    dispatch(appActions.getSprPrava())
    dispatch(appActions.getSprVisit())
    dispatch(appActions.getParams())

    console.log("Initialisation")
    return null
}

export default connect((_ => ({
})))(Initialisation)