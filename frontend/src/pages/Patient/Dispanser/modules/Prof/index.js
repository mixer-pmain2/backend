import React, {useEffect, useState} from "react";
import {connect} from "react-redux";

import {formatDateToInput} from "utility/string";
import Notify, {notifyType} from "components/Notify";
import useParams from "utility/app";
import {profVisit} from "consts/visit";
import * as apiPatient from "api/patient";
import {getHistoryVisits} from "store/actions/patient";
import {loadingAdd, loadingRemove} from "../../../../../store/actions/application";

const Prof = ({dispatch, user, application}) => {
    const [state, setState] = useState({
        error: ""
    })
    const params = useParams(application.params)
    const [dateRange, setDateRange] = useState({
        min: "",
        max: ""
    })
    const [form, setForm] = useState({
        uch: user?.section[user.unit]?.[0] || "",
        unit: user.unit,
        home: false,
        date: formatDateToInput(new Date()),
        dockId: user.id,
        count: 0
    })


    const notifySuccess = (message) => Notify(notifyType.SUCCESS, message)()
    const notifyError = (message) => Notify(notifyType.WARNING, message)()

    const onChangeHome = (e) => {
        setForm({
            ...form,
            home: e.target.checked
        })
    }

    const setDate = (e) => {
        setForm({
            ...form,
            date: e.target.value
        })
    }

    const handleReset = () => {
        setForm({
            ...form,
            count: 0
        })
    }

    const submitProf = (e) => {
        e.preventDefault()
        console.log(form)
        apiPatient.newProf(form)
            .then(res => {
                if (res?.success) {
                    const loaderName = "upd_visit"
                    notifySuccess("Проф. осомтры записаны")
                    handleReset()
                    dispatch(loadingAdd(loaderName))
                    dispatch(getHistoryVisits({id: 306258, cache: false}))
                        .finally(() => {
                            dispatch(loadingRemove(loaderName))
                        })
                    setState({
                        ...state,
                        error: ""
                    })
                }
                if (res?.error) {
                    notifyError(res?.message)
                    setState({
                        ...state,
                        error: res?.message
                    })
                }
            })
    }

    const handleOnChangeForm = (e) => {
        setForm({
            ...form,
            [e.target.name]: e.target.value
        })
    }

    const handleOnChangeCountVisit = (e) => {
        let value = e.target.value;
        if (value > profVisit.maxCount) value = profVisit.maxCount
        setForm({
            ...form,
            count: Number(value)
        })
    }

    const ViewUch = ({data}) => {

        return <div className="mb-3">
            <label className="form-label">Участок</label>
            <div>
                <select className="form-select form-select-sm" value={form.uch} onChange={handleOnChangeForm}>
                    {data?.map((v, i) =>
                        <option value={v} key={i}>{v}</option>)
                    }
                </select>
            </div>
        </div>
    }

    useEffect(() => {
        let action = false

        if (action === false) {
            setDateRange({
                min: params.visit.minDate,
                max: params.visit.maxDate
            })
        }
        return () => {
            action = true
        }
    }, [])

    return <div>
        <div className="d-flex align-items-center">
            <div className="form-check form-switch" style={{marginRight: 15}}>
                <input className="form-check-input" type="checkbox" id="home" onChange={onChangeHome}
                       checked={form.home}/>
                <label className="form-check-label" htmlFor="home">Вне диспансера</label>
            </div>
            <div className="">
                <input className="form-control" style={{width: 150, marginRight: 15}} type="date"
                       value={form.date || formatDateToInput(new Date())} onChange={setDate}
                       min={dateRange.min} max={dateRange.max}/>
            </div>
        </div>
        <form onSubmit={submitProf}>
            <hr/>
            <div className="d-flex flex-column">
                <div className="d-flex flex-row">
                    <div className="d-flex flex-row align-items-center" style={{marginRight: 15}}>
                        <label className="form-label" htmlFor="count" style={{marginRight: 15}}>
                            <span>Количество</span>
                        </label>
                        <input
                            className="form-control"
                            style={{width: 70}}
                            type="number"
                            min={0}
                            max={profVisit.maxCount}
                            value={form.count}
                            onChange={handleOnChangeCountVisit}
                            name={"count"}
                            id={"count"}
                        />
                    </div>
                </div>
                <div className="d-flex flex-column justify-content-start" style={{width: 200}}>
                    <ViewUch data={user.section[user.unit]}/>
                    <span className="text-danger">{state.error}</span>
                    <input className="btn btn-primary" type="submit" value={"Записать профосмотры"}/>
                </div>
            </div>
        </form>

    </div>
}

export default connect(state => ({
    user: state.user,
    application: state.application
}))(Prof)