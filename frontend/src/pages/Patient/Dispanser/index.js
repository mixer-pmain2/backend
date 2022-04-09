import React, {useState} from "react";
import Menu from "./Menu/Menu";
import Visit from "./modules/Visit";
import History from "./modules/History";
import Custody from "./modules/Custody";
import GroupWork from "./modules/GroupWork";
import Invalid from "./modules/Invalid";
import Uchet from "./modules/Uchet";
import Vaccinations from "./modules/Vaccinations";
import Infection from "./modules/Infection";
import {dispanserSubModules} from "../../../consts/app";

const AmbTabs = [
    dispanserSubModules.visit, dispanserSubModules.history, dispanserSubModules.uchet,
    dispanserSubModules.invalid, dispanserSubModules.custody, dispanserSubModules.groupWork,
    dispanserSubModules.vaccinations, dispanserSubModules.infection
]

const Dispanser = (props) => {
    const {patient} = props
    const [currentTab, setCurrentTab] = useState(AmbTabs[0].id)

    return <div>
        {patient?.id && <Menu
            tabs={AmbTabs}
            curTab={currentTab}
            onChange={setCurrentTab}
            style={{marginBottom: 20}}
        />}
        {currentTab === dispanserSubModules.visit.id && <Visit/>}
        {currentTab === dispanserSubModules.history.id && <History {...props}/>}
        {currentTab === dispanserSubModules.uchet.id && <Uchet {...props}/>}
        {currentTab === dispanserSubModules.invalid.id && <Invalid {...props}/>}
        {currentTab === dispanserSubModules.custody.id && <Custody {...props}/>}
        {currentTab === dispanserSubModules.groupWork.id && <GroupWork {...props}/>}
        {currentTab === dispanserSubModules.vaccinations.id && <Vaccinations {...props}/>}
        {currentTab === dispanserSubModules.infection.id && <Infection {...props}/>}
    </div>
}

export default Dispanser