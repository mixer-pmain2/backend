import React, {useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import {connect} from "react-redux"

import * as patientActions from "../../store/actions/patient"
import * as appActions from "../../store/actions/application"

import Layout from "../Layout";

import {linkDict} from "../../routes";

import {formatDate} from "../../utility/string";
import Dispanser from "./Dispanser";
import Notify, {notifyType} from "../../components/Notify";

const PatientDetail = ({onReset, patient}) => {

    const patItem = (text) =>
        <span className="input-group-text bg-light" id="basic-addon3" style={{marginRight: 10}}>
            {text}
        </span>

    const patLabel = (title, text) => <div className="input-group w-auto" style={{marginRight: 10}}>
        <span className="input-group-text">{title}</span>
        <span className="input-group-text bg-light" id="basic-addon3">
            {text}
        </span>
    </div>

    const lastUchet = patient?.uchet?.length ? patient?.uchet[0] : {}

    return <div>
        <div className="mb-3 d-flex flex-row justify-content-between">
            <a href="#" className="btn btn-outline-secondary" style={{marginRight: 10}}
               onClick={onReset}>Сброс</a>
            <div className="input-group" style={{marginRight: 10, width: 300}}>
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
        <div className="mb-3 d-flex flex-row">
            {patLabel("Участок", lastUchet?.section)}
            {patLabel("Учет", lastUchet?.categoryS)}
            {patLabel("Диагноз учета", lastUchet?.diagnose)}
            {patLabel("Адрес", patient?.address)}
        </div>
    </div>
}

const GetPatient = (props) => {
    const {dispatch, patient, application} = props

    const navigation = useNavigate()
    const {id} = useParams()

    const foundById = (id) => {
        return dispatch(patientActions.findById({id}))
            .then(res => {
                dispatch(patientActions.select(res))
                dispatch(patientActions.getUchet({id: patient.id}))
                return res
            })
    }

    const onReset = () => {
        dispatch(patientActions.reset())
        navigation(linkDict.findPatient)
    }

    const onEscDown = (e) => {
        if (e.code === 'Escape') onReset()
    }

    useEffect(() => {
        let action = false

        if (!action) {
            if (patient?.id) {
                dispatch(appActions.enableLoading())
                foundById(id)
                dispatch(appActions.disableLoading())
            } else {
                navigation(linkDict.findPatient)
            }

            document.addEventListener("keydown", onEscDown, false);
        }

        return () => {
            action = true
            document.removeEventListener("keydown", onEscDown, false);
        }
    }, [])

    // useEffect(() => {
        // dispatch(appActions.enableLoading())
        // if (patient?.id && !patient?.uchet) {
        //     dispatch(patientActions.getUchet({id: patient.id}))
        // }
        // dispatch(appActions.disableLoading())
    // }, [patient?.id])

    return <Layout>
        <PatientDetail
            patient={patient}
            onReset={onReset}
        />
        <hr style={{marginTop: 10, marginBottom: 25}}/>
        <Dispanser {...props}/>
    </Layout>
}

export default connect(state => ({
    patient: state.patient,
    application: state.application
}))(GetPatient)