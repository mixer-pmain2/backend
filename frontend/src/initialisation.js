import {useEffect} from "react";
import {connect} from "react-redux";
import * as appActions from "./store/actions/application"


// Не добавлять состояния иначе будет часто вызываться
const Initialisation = ({dispatch, user}) => {


    console.log("Initialisation")
    return null
}

export default connect((_ => ({
})))(Initialisation)