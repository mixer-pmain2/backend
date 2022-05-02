import React, {useEffect, useState} from "react";
import {connect} from "react-redux";

import Table from "components/Table";
import * as patientActions from "store/actions/patient";
import {formatDate, formatDateToInput} from "utility/string";
import useParams from "utility/app";
import Patient from "classes/Patient";




const HistoryUchet = ({patient}) => {
    const sindrom = [
        {title: "Sim1", value: "Sim1T"},
        {title: "Sim2", value: "Sim2T"},
        {title: "Sim3", value: "Sim3T"},
        {title: "Sim4", value: "Sim4T"},
    ]
    const somat = [
        {title: "Som1", value: "Som1T"},
        {title: "Som2", value: "Som2T"},
        {title: "Som3", value: "Som3T"},
        {title: "Som4", value: "Som4T"},
    ]


    const mapper = (row) => {
        return <>
            <td>{row.section}</td>
            <td>{formatDate(row.date)}</td>
            <td>{row.categoryS}</td>
            <td>{row.reasonS}</td>
            <td>{row.diagnose}</td>
            <td>{row.diagnoseS}</td>
            <td>{row.dockName}</td>
        </>
    }

    const diagItems = (title, value) =>
        <div key={title}>
            <span>{title}</span> <span>{value}</span>
        </div>

    return <div>
        <Table
            columns={["Участок", "Дата", "Категори", "Причина", "Учет", "Расшифровка диагноза", "Врач"]}
            data={patient?.uchet || []}
            mapper={mapper}
        />
        <div className="mb-4">
            <h6>Синдромы</h6>
            <div style={{marginLeft: 15}}>
                {sindrom.map((v, i) => diagItems(v.title, v.value))}

            </div>
        </div>
        <div>
            <h6>Хронические соматические заболевания</h6>
            <div style={{marginLeft: 15}}>
                {somat.map((v, i) => diagItems(v.title, v.value))}
            </div>
        </div>
    </div>
}

const NewUchet = ({onClose}) => {
    return <div>
        newuchet
        <button className="btn btn-outline-danger" onClick={onClose}>Закрыть</button>
    </div>
}

const NewSomSin = ({onClose}) => {
    return <div>
        newsomsin
        <button className="btn btn-outline-danger" onClick={onClose}>Закрыть</button>
    </div>
}

const T_UCHET = 0
const T_NEW = 1
const T_SIN_SOM = 2

const Uchet = ({dispatch, patient, application}) => {
    const [state, setState] = useState({
        tab: T_UCHET,
        enableTransfer: false
    })
    const [form, setForm] = useState({
        date: formatDateToInput(new Date()),
    })
    const [dateRange, setDateRange] = useState({
        min: "",
        max: ""
    })
    const params = useParams(application.params)

    const onChangeForm = (e) => {
        setForm({
            ...form,
            [e.target.name]: e.target.value
        })
    }

    const handleNew = () => setState({...state, tab: T_NEW})
    const handleSomSin = () => setState({...state, tab: T_SIN_SOM})
    const closeTab = () => setState({...state, tab: T_UCHET})

    useEffect(() => {
        dispatch(patientActions.getUchet({id: patient.id}))
        setDateRange({
            min: params.visit.minDate,
            max: params.visit.maxDate
        })
        const p = new Patient(patient)
        if (p.getLastUchet().category) setState({...state, enableTransfer: true})
    }, [])

    return <div>
        <div className="d-flex flex-row justify-content-between">
            <div className="d-flex flex-row align-items-center">
                <button className="btn btn-outline-primary" style={{marginRight: 15}} onClick={handleNew}>+</button>
                {state.enableTransfer && <button className="btn btn-outline-primary" style={{marginRight: 15}}>
                    Прием с других участков
                </button>}
                <input className="form-control" style={{width: 150, marginRight: 15}} type="date" name="date"
                       value={form.date || formatDateToInput(new Date())} onChange={onChangeForm}
                       min={dateRange.min} max={dateRange.max}/>
                <span></span>
            </div>
            <button className="btn btn-outline-primary" onClick={handleSomSin}>Синдром и хронические заболевания</button>
        </div>
        <hr/>
        {state.tab === T_UCHET && <HistoryUchet patient={patient}/>}
        {state.tab === T_NEW && <NewUchet onClose={closeTab}/>}
        {state.tab === T_SIN_SOM && <NewSomSin onClose={closeTab}/>}
    </div>
}

export default connect(state => ({
    patient: state.patient,
    application: state.application
}))(Uchet)