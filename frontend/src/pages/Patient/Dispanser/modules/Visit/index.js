import React, {useState} from "react";
import {connect} from "react-redux";

import {visit_home_type_group, visit_type_group} from "../../../../../consts/visit";


const Visit = ({dispatch, application, patient, user}) => {
  const [form, setForm] = useState({
    uch: user?.section[user.unit]?.[0] || "",
    visit: 0,
    unit: user.unit,
    home: false
  })
  const sprVisit = application?.spr?.visit || {};

  const isValidVisit = () => {
    if (form.visit === 0) return false
    if (form.uch === "") return false

    return true
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

  const SideUch = ({data}) => {

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
      uch: e.target.value
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

  const submitNewVisit = (e) => {
    e.preventDefault()
    console.log(e)
    console.log(form)
  }
  console.log(form)
  return <div>
    <div className="form-check form-switch">
      <input className="form-check-input" type="checkbox" id="home" onChange={onChangeHome} checked={form.home}/>
      <label className="form-check-label" htmlFor="home">На дому</label>
    </div>
    <form onSubmit={submitNewVisit}>
      <div className="d-flex flex-row">
        <div className="d-flex flex-row, flex-wrap">
          {(form.home ? visit_home_type_group : visit_type_group).map((v, i) =>
            <GroupVisit group={v} key={i}/>
          )}
        </div>
        <div className="d-flex flex-column justify-content-start" style={{width: 200}}>
          <SideUch data={user.section[user.unit]}/>
          <input type="submit" className="btn btn-primary" value="Записать" disabled={!isValidVisit()}/>
        </div>

      </div>
    </form>
  </div>
}

export default connect(state => ({
  application: state.application,
  patient: state.application,
  user: state.user
}))(Visit)