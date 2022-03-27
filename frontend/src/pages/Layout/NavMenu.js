import React from "react";
import {Link} from "react-router-dom";

import {linkDict} from "../../routes";

import {capitalizeFirstLetter} from "../../utility/string";


const NavMenu = ({onLogout, user, app, patient}) => {
    const urlPatient = () =>
        linkDict.patient.replaceAll(':id', patient?.id)

    const fio = capitalizeFirstLetter(user?.lname.toLowerCase())+" "
        +capitalizeFirstLetter(user?.fname.toLowerCase())[0]+"."
        +capitalizeFirstLetter(user?.sname.toLowerCase())[0]+"."

    return <nav className="navbar navbar-expand-lg navbar-light bg-light">
        <div className="container-fluid">
            <Link to={linkDict.start} className="navbar-brand">КПБ</Link>
            <div className="collapse navbar-collapse" id="navbarSupportedContent">
                <ul className="navbar-nav me-auto mb-2 mb-lg-0">
                    {patient?.id && <li className="nav-item">
                        <Link to={urlPatient()} className="nav-link">{app?.spr?.unit[user?.unit]}</Link>
                    </li>}
                    <li className="nav-item">
                        <Link to={linkDict.findPatient} className="nav-link">Администрирование</Link>
                    </li>
                </ul>
            </div>
            <div>
                <button className="btn btn-danger">
                    <Link to={linkDict.admin} className="nav-link link-light" style={{padding: 0}}>admin</Link>
                </button>
                <span className="m-2 navbar-text">Пользователь: {fio}</span>
                <button className="btn btn-primary" onClick={onLogout}>Выход</button>
            </div>
        </div>
    </nav>
}
export default NavMenu