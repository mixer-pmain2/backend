import React, {useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import {connect} from "react-redux"

import * as patientActions from "../../store/actions/patient"
import * as appActions from "../../store/actions/application"

import Menu, {AmbTabs} from "./Ambulance/Menu/Menu";
import Layout from "../Layout";

import {linkDict} from "../../routes";

import {formatDate} from "../../utility/string";

const PatientDetail = ({onReset, patient}) => {

    const patItem = (text) =>
        <span className="input-group-text bg-light" id="basic-addon3" style={{marginRight: 10}}>
            {text}
        </span>

    return <form onSubmit={e => e.preventDefault()} autoComplete="off">
        <div className="mb-3 d-flex flex-row">
            <a href="#" className="btn btn-outline-secondary" style={{marginRight: 10}}
               onClick={onReset}>Сброс</a>
            <div className="input-group w-25" style={{marginRight: 10}}>
                <span className="input-group-text">Шифр</span>
                <input
                    type="text"
                    className="form-control bg-light"
                    autoFocus={true}
                    id="patientId"
                    name="patientId"
                    maxLength={10}
                    value={patient.id ? Number(patient?.id) : ""}
                    readOnly={true}
                />
            </div>
            <div className="input-group w-100" style={{marginRight: 10}}>
                <span className="input-group-text" id="basic-addon3">Ф.И.О.</span>
                <input
                    type="text"
                    className="form-control bg-light"
                    autoFocus={true}
                    value={(patient?.lname || "") + " " + (patient?.fname || "") + " " + (patient?.sname || "")}
                    readOnly={true}
                />
            </div>

            {patItem(patient?.sex)}
            {patItem(formatDate(patient?.bday))}
            {patItem(patient?.snils)}
            <button className="input-group-text btn btn-outline-primary">Найти</button>
        </div>
    </form>
}

const GetPatient = ({dispatch, patient, application}) => {
    const [patientData, setPatientData] = useState({})
    const [currentTab, setCurrentTab] = useState(AmbTabs[0])

    const navigation = useNavigate()
    const {id} = useParams()

    const foundById = (id) => {
        dispatch(patientActions.findById({id}))
            .then(res => {
                dispatch(appActions.enableLoading())
                dispatch(patientActions.select(res))
            })
            .finally(_ => {
                dispatch(appActions.disableLoading())
            })
    }

    const onReset = () => {
        navigation(linkDict.dispFindPatient)
    }

    const onEscDown = (e) => {
        if (e.code === 'Escape') onReset()
    }

    useEffect(() => {
        let action = false

        if (!action) {
            if (Number(id) === patient?.id) {
                setPatientData(patient)
            } else foundById(id)

            document.addEventListener("keydown", onEscDown, false);
        }

        return () => {
            action = true
            document.removeEventListener("keydown", onEscDown, false);
        }
    }, [])

    return <Layout>
        <PatientDetail
            patient={patientData}
            onReset={onReset}
        />
        {patient?.id && <Menu curTab={currentTab} onChange={setCurrentTab}/>}
    </Layout>
}

export default connect(state => ({
    patient: state.patient,
    application: state.application
}))(GetPatient)