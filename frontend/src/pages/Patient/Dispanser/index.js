import React, {useState} from "react";
import {connect} from "react-redux";

import Menu from "./Menu/Menu";
import Visit from "./modules/Visit";
import History from "./modules/History";
import Custody from "./modules/Custody";
import GroupWork from "./modules/GroupWork";
import Invalid from "./modules/Invalid";
import Uchet from "./modules/Uchet";
import Vaccinations from "./modules/Vaccinations";
import Infection from "./modules/Infection";
import UKL from "./modules/UKL";
import Prof from "./modules/Prof";

import {dispanserSubModules} from "consts/app";


const AmbTabs = [
    dispanserSubModules.history, dispanserSubModules.visit, dispanserSubModules.prof, dispanserSubModules.uchet,
    dispanserSubModules.invalid, dispanserSubModules.custody, dispanserSubModules.groupWork,
    dispanserSubModules.vaccinations, dispanserSubModules.infection, dispanserSubModules.UKL
]

const Dispanser = (props) => {
    const {patient, user} = props
    const [currentTab, setCurrentTab] = useState(-1)

    const tabsFilter = (tabs) => {
        return tabs.filter(v => (v.unit & user.unit) && (v.access & user.access[user.unit]))
    }

    return <div>
        {patient?.id && <Menu
            tabs={tabsFilter(AmbTabs)}
            curTab={currentTab}
            onChange={setCurrentTab}
            style={{marginBottom: 20}}
        />}
        {currentTab === dispanserSubModules.visit.id && <Visit/>}
        {currentTab === dispanserSubModules.prof.id && <Prof/>}
        {currentTab === dispanserSubModules.history.id && <History/>}
        {currentTab === dispanserSubModules.uchet.id && <Uchet/>}
        {currentTab === dispanserSubModules.invalid.id && <Invalid/>}
        {currentTab === dispanserSubModules.custody.id && <Custody/>}
        {currentTab === dispanserSubModules.groupWork.id && <GroupWork/>}
        {currentTab === dispanserSubModules.vaccinations.id && <Vaccinations/>}
        {currentTab === dispanserSubModules.infection.id && <Infection/>}
        {currentTab === dispanserSubModules.UKL.id && <UKL/>}
    </div>
}

export default connect(state => ({
    patient: state.patient,
    user: state.user
}))(Dispanser)