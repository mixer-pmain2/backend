import React, {useCallback, useEffect} from "react"
import {connect} from "react-redux";

import Layout from "../Layout";

import Card from "./components/Card";

import * as userActions from "../../store/actions/user";
import {linkDict} from "../../routes";
import {useNavigate} from "react-router-dom";
import * as appActions from "../../store/actions/application";

const podrList = [
    {
        title: "Взрослая психиатрия",
        img: require("../../assets/images/295067-3f51b5.svg"),
        podr: 1
    },
    {
        title: "Психотерапия",
        img: require("../../assets/images/295067-3f51b5.svg"),
        podr: 2
    },
    {
        title: "Специалисты",
        img: require("../../assets/images/295067-3f51b5.svg"),
        podr: 1024
    },
    {
        title: "Суицидология",
        img: require("../../assets/images/295067-3f51b5.svg"),
        podr: 4
    },
    {
        title: "ОИЛС",
        img: require("../../assets/images/295067-3f51b5.svg"),
        podr: 8
    },
    {
        title: "Детский диспансер",
        img: require("../../assets/images/2025726-3f51b5.svg"),
        podr: 16
    },
    {
        title: "Подростковая психиатрия",
        img: require("../../assets/images/295067-3f51b5.svg"),
        podr: 16777216
    },
    {
        title: "Детская консультация",
        img: require("../../assets/images/295067-3f51b5.svg"),
        podr: 33554432
    }
]

const StartPage = ({dispatch, user, application}) => {
    const navigate = useNavigate()

    const initUser = useCallback(() => {
        dispatch(userActions.getPrava({id: user?.id}))
        dispatch(appActions.getSprPodr())
    }, [])

    useEffect(() => {
        initUser()
    }, [])

    const handleSelectPodr = (podr) => {
        dispatch(userActions.setCurrentPodr(podr))
        navigate(linkDict.findPatient)
    }

    return <Layout>
        <div className="row">
            {podrList.map((v, i) => <div key={i} className="col-sm-3 mb-4">
                <Card title={v.title} img={v.img} onClick={_ => handleSelectPodr(v.podr)}/>
            </div>)}

        </div>
    </Layout>
}

export default connect((state) => ({
    user: state.user,
    application: state.application
}))(StartPage)