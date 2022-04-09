import React, {useCallback, useEffect} from "react"
import {connect} from "react-redux";

import Layout from "../Layout";

import Card from "./components/Card";

import * as userActions from "../../store/actions/user";
import {linkDict} from "../../routes";
import {useNavigate} from "react-router-dom";


const unitList = [
    {
        title: "Взрослая психиатрия",
        img: require("../../assets/images/295067-3f51b5.svg"),
        unit: 1
    },
    {
        title: "Психотерапия",
        img: require("../../assets/images/295067-3f51b5.svg"),
        unit: 2
    },
    {
        title: "Специалисты",
        img: require("../../assets/images/295067-3f51b5.svg"),
        unit: 1024
    },
    {
        title: "Суицидология",
        img: require("../../assets/images/295067-3f51b5.svg"),
        unit: 4
    },
    {
        title: "ОИЛС",
        img: require("../../assets/images/295067-3f51b5.svg"),
        unit: 8
    },
    {
        title: "Детский диспансер",
        img: require("../../assets/images/2025726-3f51b5.svg"),
        unit: 16
    },
    {
        title: "Подростковая психиатрия",
        img: require("../../assets/images/295067-3f51b5.svg"),
        unit: 16777216
    },
    {
        title: "Детская консультация",
        img: require("../../assets/images/295067-3f51b5.svg"),
        unit: 33554432
    }
]

const StartPage = ({dispatch, user, application}) => {
    const navigate = useNavigate()

    const initUser = useCallback(() => {
        if (user?.id) {
            dispatch(userActions.getPrava({id: user?.id}))
            dispatch(userActions.getUch({id: user?.id}))
        }
    }, [user?.id])

    useEffect(() => {
        initUser()
    }, [user?.id])

    const handleSelectUnit = (unit) => {
        dispatch(userActions.setCurrentPodr(unit))
        navigate(linkDict.findPatient)
    }

    return <Layout>
        <div className="row d-flex justify-content-center">
            {unitList.map((v, i) =>
                <Card key={i} title={v.title} img={v.img} onClick={_ => handleSelectUnit(v.unit)}/>
            )}

        </div>
    </Layout>
}

export default connect((state) => ({
    user: state.user,
    application: state.application
}))(StartPage)