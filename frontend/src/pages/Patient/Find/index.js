import React, {useEffect, useState} from "react";
import {connect} from "react-redux";
import {useNavigate} from "react-router-dom";

import * as patientActions from "../../../store/actions/patient"
import * as appActions from "../../../store/actions/application"

import {app} from "../../../consts/app"
import {linkDict} from "../../../routes";
import {formatDate} from "../../../utility/string";

import Layout from "../../Layout";

import Table from "../../../components/Table";
import notify, {notifyType} from "../../../components/Notify";


const FindPatient = ({dispatch, application, patient, type = app.DISPANSER}) => {
    const [isFounded, setIsFounded] = useState(false)
    const [fio, setFio] = useState("")
    const [patientId, setPatientId] = useState("")
    const [foundPatient, setFoundPatient] = useState([])

    const navigate = useNavigate()

    const upper = (text) => text.toUpperCase()

    const onReset = () => {
        setIsFounded(false)
        setFio("")
        setPatientId("")
    }

    const notifyInfo = (message) => notify(type = notifyType.INFO, message)()


    const findByFio = fio => {
        dispatch(appActions.enableLoading())
        dispatch(patientActions.findByFio({fio}))
            .then(res => {
                setFoundPatient(res)
            })
            .finally(() => {
                dispatch(appActions.disableLoading())
                setIsFounded(true)
            })
    }

    const findById = id => {
        dispatch(appActions.enableLoading())
        dispatch(patientActions.findById({id}))
            .then(res => {
                if (res?.id) {
                    dispatch(patientActions.select(res))
                    navigate(linkDict.patient.replaceAll(":id", id))
                } else {
                    notifyInfo("Нет данных")
                }
            })
            .finally(_ => {
                dispatch(appActions.disableLoading())
            })
    }

    const handleFindPatient = (e) => {
        e.preventDefault()
        const [patientId, fio] = e.target
        // if (fio.value.length < 3 && fio.value.length > 0) return notifyInfo("Нужно больше информации")

        if (patientId.value.length > 0) {
            findById(patientId.value)
        } else if (fio.value.length > 0) {
            findByFio(fio.value)
        } else {
            notifyInfo("Нужно больше информации")
        }
    }

    const handleSelectPatient = (e) => {
        dispatch(patientActions.select(e))
        navigate(linkDict.patient.replaceAll(":id", e.id))
    }

    const mapper = (row) => {
        return <>
            <th scope="row">
                {row.id}
            </th>
            <td>
                {row.lname}
            </td>
            <td>
                {row.fname}
            </td>
            <td>
                {row.sname}
            </td>
            <td>
                {formatDate(row.bday)}
            </td>
        </>
    }

    const onEscDown = (e) => {
        if (e.code === 'Escape') onReset()
    }

    useEffect(() => {
        let action = false

        if (!action) {
            document.addEventListener("keydown", onEscDown, false);
        }

        return () => {
            action = true
            document.removeEventListener("keydown", onEscDown, false);
        }
    }, [])

    return <Layout>
        <form onSubmit={handleFindPatient} autoComplete="off">
            <div className="mb-3 d-flex flex-row">
                <a href="#" className="btn btn-outline-secondary" style={{marginRight: 10}} onClick={onReset}>Сброс</a>
                <div className="input-group w-25" style={{marginRight: 10}} >
                    <span className="input-group-text">Шифр</span>
                    <input
                        type="text"
                        className="form-control"
                        id="patientId"
                        name="patientId"
                        autoFocus={true}
                        maxLength={10}
                        onChange={e => {
                            e.target.value.length > 0 && setPatientId(e.target.value)
                            setFio("")
                        }}
                        value={Number(patientId) || ""}
                    />

                </div>
                <div className="input-group w-100" style={{marginRight: 10}}>
                    <span className="input-group-text">Ф.И.О.</span>
                    <input
                        type="text"
                        className="form-control"
                        id="fio"
                        name="fio"
                        maxLength={255}
                        onChange={e => {
                            setFio(e.target.value)
                            setPatientId("")
                        }}
                        value={upper(fio)}
                    />
                </div>
                <button className="input-group-text btn btn-outline-primary">Найти</button>
            </div>
        </form>
        {isFounded && <Table
            columns={["Шифр", "Фамилия", "Имя", "Отчество", "Дата рождения"]}
            data={foundPatient}
            mapper={mapper}
            onDoubleClick={handleSelectPatient}
            selecting={true}
        />}
    </Layout>
}

export default connect(state => ({
    patient: state.patient,
    application: state.application
}))(FindPatient)