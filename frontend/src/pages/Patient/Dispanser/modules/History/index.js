import React, {useEffect, useState} from "react";
import {connect} from "react-redux";

import * as patientActions from "../../../../../store/actions/patient"
import * as appActions from "../../../../../store/actions/application"

import Table from "../../../../../components/Table";
import {formatDate, shorty} from "../../../../../utility/string";


const LOADING_VISIT = "history_visit"
const LOADING_HOSPITAL = "history_hospital"


const HistoryVisit = ({patient}) => {
    const [selectRow, setSelectRow] = useState(patient?.visit?.[0])
    const mapper = (row) => {
        return <>
            <td>
                {formatDate(row.date)}
            </td>
            <td>
                {row.dockName}
            </td>
            <td>
                {row.diag}
            </td>
            <td>
                {row.reason}
            </td>
        </>
    }

    return <div>
        <span>Диагноз приема </span><span title={selectRow?.diagS}>{shorty(selectRow?.diagS, 115)}</span>
        <Table
            columns={["Дата", "Врач", "Диагноз", "Причина"]}
            data={patient.visit || []}
            mapper={mapper}
            selecting={true}
            onClick={setSelectRow}
        />
    </div>
}

const HistoryHospital = ({patient}) => {
    const [selectRow, setSelectRow] = useState(patient?.visit?.[0])
    const mapper = (row) => {
        return <>
            <td>
                {formatDate(row.dateStart)}
            </td>
            <td>
                {formatDate(row.dateEnd)}
            </td>
            <td>
                {row.otd}
            </td>
            <td>
                {row.diagStart}
            </td>
            <td>
                {row.diagEnd}
            </td>
        </>
    }

    return <div>
        <span>Диагноз выписки </span><span title={selectRow?.diagEndS}>{shorty(selectRow?.diagEndS)}</span>
        <Table
            columns={["Дата поступления", "Дата выписки", "Отделение", "Диагноз поступления", "Диагноз выписки"]}
            data={patient.hospital || []}
            mapper={mapper}
            selecting={true}
            onClick={setSelectRow}
        />

    </div>
}

const History = ({dispatch, patient}) => {

    useEffect(() => {
        if (!patient?.visit) {
            dispatch(appActions.loadingAdd(LOADING_VISIT))
            dispatch(patientActions.getHistoryVisits({id: patient.id}))
                .finally(_ => {
                    dispatch(appActions.loadingRemove(LOADING_VISIT))
                })
        }
        if (!patient?.hospital) {
            dispatch(appActions.loadingAdd(LOADING_HOSPITAL))
            dispatch(patientActions.getHistoryHospital({id: patient.id}))
                .finally(_ => {
                    dispatch(appActions.loadingRemove(LOADING_HOSPITAL))
                })
        }

    }, [])

    useEffect(() => {
        console.log(patient)
    })

    return <div>
        <div style={{marginBottom: 20}}>
            <HistoryVisit patient={patient}/>
        </div>
        <div>
            <HistoryHospital patient={patient}/>
        </div>
    </div>
}

export default connect(state => ({
    patient: state.patient,
    application: state.application
}))(History)