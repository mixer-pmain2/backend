import React, {useState} from "react";
import {connect} from "react-redux";
import UserFind from "./UserFind";

const UserAccess = (props) => {
    const {application, user, dispatch} = props
    const units = application.spr.unit
    const unitsIndex = Object.keys(units)
    const [curTab, setCurTab] = useState(unitsIndex[0])

    const access = application.spr.access
    let accessFiltered = access.filter((v) => v.unit === Number(curTab))
    const [isCheckedAccess, setIsCheckedAccess] = useState(
        accessFiltered.map(v => ({[v.code]: (user?.access?.[curTab] & v.code) > 0}))
    )

    return <div>
        <UserFind/>
        <hr/>
        <div className="d-flex">
            <div className="w-25">
                <nav className="nav nav-pills flex-column">
                    {
                        unitsIndex.map((v, i) =>
                            <button className={`nav-link w-100 align-items-start ${v === curTab ? "active" : ""}`}
                                    key={i}
                                    onClick={_ => setCurTab(v)}>
                                <span>{units[v]}</span>
                            </button>
                        )
                    }

                </nav>
            </div>
            <div className="w-75" style={{marginLeft: 35}}>
                {
                    accessFiltered.map((v, i) => {
                        return <div key={i} className="form-check form-switch">
                            <input className="form-check-input" type="checkbox" id={`flexSwitchCheckDefault${v.code}`}
                                   checked={isCheckedAccess[v.code]}
                                   onChange={_ => setIsCheckedAccess({
                                       ...isCheckedAccess,
                                       [v.code]: !isCheckedAccess[v.code]
                                   })}
                            />
                            <label className="form-check-label" htmlFor={`flexSwitchCheckDefault${v.code}`}>
                                {v.name}
                            </label>
                        </div>
                    })
                }

            </div>
        </div>
    </div>
}

export default connect(state => ({
    user: state.user,
    application: state.application
}))(UserAccess)