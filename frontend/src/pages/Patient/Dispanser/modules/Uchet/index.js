import React from "react";
import {connect} from "react-redux";

const HistoryUchet = () => {

}

const Uchet = () => {
    return <div>
        <HistoryUchet/>
    </div>
}

export default connect(state => ({
    patient: state.patient
}))(Uchet)