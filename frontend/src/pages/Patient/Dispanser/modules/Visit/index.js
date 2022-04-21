import React, {useEffect, useState} from "react";
import {connect} from "react-redux";

import DiagnoseTree, {getTypeDiagsModal} from "pages/Patient/components/DiagnoseTree";
import Modal, {BTN_NO, BTN_YES} from "components/Modal";
import Notify, {notifyType} from "components/Notify";

import {getHistoryVisits} from "store/actions/patient";
import * as apiPatient from "api/patient"

import {getAge, stringToDate} from "utility/date";
import {formatDateToInput} from "utility/string";
import useParams from "utility/app"

import {Unit} from "consts/user";
import {messageNotValidAge, srcReason, visitHomeTypeGroup, visitTypeGroup} from "consts/visit";


const Visit = ({dispatch, application, patient, user}) => {
  const [state, setState] = useState({
    error: "",
    src: -1,
    isAcceptNotValidAge: false,
    isOpenConfirmAge: false,
    isOpenDiagModal: false,
    isSubmit: false,
  })
  const params = useParams(application.params)
  const srcEnable = user.unit === 16 || user.unit === 16777216 || user.unit === 33554432 ? true : false
  const [dateRange, setDateRange] = useState({
    min: "",
    max: ""
  })
  const [form, setForm] = useState({
    uch: user?.section[user.unit]?.[0] || "",
    visit: 0,
    unit: user.unit,
    home: false,
    diagnose: "",
    date: formatDateToInput(new Date()),
    patientId: patient.id,
    patientBDay: patient.bday,
    dockId: user.id,
    src: -1
  })
  const sprVisit = application?.spr?.visit || {};
  const messageNotValidYearOld = user.unit === 1 ? messageNotValidAge.sup : messageNotValidAge.sub

  const notifySuccess = (message) => Notify(notifyType.SUCCESS, message)()
  const notifyError = (message) => Notify(notifyType.WARNING, message)()

  const isValidVisit = () => {
    // if (form.visit === 0) return false
    if (form.uch === "") return false
    // if (form.diagnose === "") return false

    return true
  }

  const getLastDiag = () => {
    return patient?.visit?.length === 0 || user.unit === 1024 ? "" : patient?.visit?.[0]?.diag
  }

  const ItemVisit = ({value}) => <label className="mb-1">
    <input type="checkbox" className="form-check-input fs-6" value={value} onChange={onClickVisit}
           checked={(form.visit & value) > 0}/>
    <span className="text p-2 fs-6">{sprVisit[value]}</span>
  </label>

  const GroupVisit = ({group}) =>
    <div className="mb-4 d-flex flex-column" style={{marginRight: 30, width: 500}}>
      <hr/>
      {group.map((v, i) =>
        <ItemVisit value={v} key={i}/>
      )}
    </div>

  const ViewUch = ({data}) => {

    return <div className="mb-3">
      <span className="fs-5">Участок</span>
      <div>
        <select className="form-select form-select-sm" value={form.uch} onChange={onChangeUch}>
          {data?.map((v, i) =>
            <option value={v} key={i}>{v}</option>)
          }
        </select>
      </div>
    </div>
  }

  const onChangeUch = (e) => {
    setForm({
      ...form,
      uch: Number(e.target.value)
    })
  }

  const onClickVisit = (e) => {
    const isChecked = e.target.checked
    const value = Number(e.target.value)
    setForm({
      ...form,
      visit: isChecked ? form.visit + value : form.visit - value
    })
  }

  const onChangeHome = (e) => {
    setForm({
      ...form,
      home: e.target.checked
    })
  }

  const isValidAge = () => {
    const age = getAge(stringToDate(patient.bday), stringToDate(form.date))
    if (user.unit === 1 && age < 18) return false
    if (([Unit["Детский диспансер"], Unit["Детская консультация"], Unit["Подростковая психиатрия"]].indexOf(user.unit) + 1)
      && age > 17) return false

    return true
  }

  const handleAcceptAgeYes = () => {
    setState({
      ...state,
      isAcceptNotValidAge: true,
      isOpenConfirmAge: false,
      isSubmit: true
    })
    // handleNewVisit()
  }

  const handleAcceptAgeNo = () => {
    setState({
      ...state,
      isAcceptNotValidAge: false,
      isOpenConfirmAge: false,
      isSubmit: false
    })
  }

  const submitNewVisit = (e) => {
    e.preventDefault()
    handleNewVisit()
  }

  const handleNewVisit = () => {
    console.log(form)
    if (form.diagnose === "" || form.diagnose.length < 3) {
      setState({
        ...state,
        isOpenDiagModal: true
      })
      return
    }
    console.log(!state.isAcceptNotValidAge, !isValidAge())
    if (!state.isAcceptNotValidAge && !isValidAge()) {
      setState({
        ...state,
        isOpenConfirmAge: true
      })
      return
    }

    apiPatient.newVisit(form)
      .then(res => {
        if (res?.success) {
          notifySuccess("Посещение записано")
          handleReset()
          dispatch(getHistoryVisits({id: patient.id, cache: false}))
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

  const handleReset = () => {
    setForm({
      ...form,
      visit: 0
    })
    setState({
      ...state,
      error: "",
      src: -1,
      isAcceptNotValidAge: false,
      isSubmit: false
    })
  }

  const setDate = (e) => {
    setForm({
      ...form,
      date: e.target.value
    })
  }

  const onSelectDiag = (diagnose) => {
    setState({
      ...state,
      isOpenDiagModal: false,
      isSubmit: true
    })
    setForm({
      ...form,
      diagnose: diagnose,
    })
  }

  const onChangeSRC = (e) => {
    let srcValue = -1
    if (e.target.checked) srcValue = e.target.value

    setState({
      ...state,
      src: srcValue
    })
    setForm({
      ...form,
      src: Number(srcValue),
    })
  }

  useEffect(() => {
    if (state.isSubmit) {
      handleNewVisit()
    }

  }, [state.isSubmit])

  useEffect(() => {
    let action = false

    if (action === false) {
      !patient.visit && dispatch(getHistoryVisits({id: patient.id}))
      setDateRange({
        min: params.visit.minDate,
        max: params.visit.maxDate
      })
    }
    return () => {
      action = true
    }
  }, [])

  useEffect(() => {
    setForm({
      ...form,
      diagnose: getLastDiag()
    })
  }, [patient.visit])
  console.log(state)
  return <div>
    <div className="d-flex align-items-center">
      <div className="form-check form-switch" style={{marginRight: 15}}>
        <input className="form-check-input" type="checkbox" id="home" onChange={onChangeHome} checked={form.home}/>
        <label className="form-check-label" htmlFor="home">На дому</label>
      </div>
      <div>
        <input className="form-control" style={{width: 150, marginRight: 15}} type="date"
               value={form.date || formatDateToInput(new Date())} onChange={setDate}
               min={dateRange.min} max={dateRange.max}/>
      </div>
      <div>
        <span>Диагноз: {form.diagnose}</span>
      </div>
    </div>
    <form onSubmit={submitNewVisit}>
      <div className="d-flex flex-row">
        <div className="d-flex flex-row, flex-wrap">
          {(form.home ? visitHomeTypeGroup : visitTypeGroup).map((v, i) =>
            <GroupVisit group={v} key={i}/>
          )}
        </div>
        <div className="d-flex flex-column justify-content-start" style={{width: 200}}>
          <ViewUch data={user.section[user.unit]}/>
          <span className="text-danger">{state.error}</span>
          <input type="submit" className="btn btn-primary" value="Записать" disabled={!isValidVisit()}/>
        </div>

      </div>
      {srcEnable && !form.home && <div>
        <h6>Направление в СРЦ</h6>
        <div className="d-flex flex-row">
          {
            srcReason.map((v, i) => <div key={i} style={{marginRight: 15}}>
              <input type="checkbox" className="form-check-input p-2" name="src_reason" id={`src_reason${i}`}
                     value={v.value} style={{marginRight: 5}} onChange={onChangeSRC} checked={state.src == v.value}/>
              <label className="form-check-label" htmlFor={`src_reason${i}`}>{v.title}</label>
            </div>)
          }
        </div>
      </div>}
    </form>
    <Modal
      style={{maxWidth: 750}}
      body={<DiagnoseTree type={getTypeDiagsModal(form.uch)} onSelect={onSelectDiag}/>}
      isOpen={state.isOpenDiagModal}
      onClose={() => setState({...state, isOpenDiagModal: false})}
    />
    <Modal
      body={<div>{messageNotValidYearOld}</div>}
      isOpen={state.isOpenConfirmAge}
      onClose={() => setState({...state, isOpenConfirmAge: false})}
      btnNum={BTN_YES + BTN_NO}
      onYes={handleAcceptAgeYes}
      onNo={handleAcceptAgeNo}

    />

  </div>
}

export default connect(state => ({
  application: state.application,
  patient: state.patient,
  user: state.user
}))(Visit)