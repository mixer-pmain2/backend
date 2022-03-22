import React from "react";
import StartPage from "./pages/Start";
import FoundPatient from "./pages/Patient/Found";
import VisitPage from "./pages/Patient/Ambulance/Visit/index";


import TestPage from "./pages/Test";
import GetPatient from "./pages/Patient";
import ProfilePage from "./pages/Profile";
import middlewareAuthRequire from "./middleware/authRequier";
import NewMiddleWare from "./middleware";
import middlewareAccessRequire from "./middleware/accessRequier";
import SignInPageState, {SignInLayoutPage} from "./pages/Signin";

export const linkDict = {
    start: "/",

    profile: "/profile",
    signin: "/signin",

    disp: "/disp",
    patient: "/patient/:id",
    findPatient: "/patient/find",
    dispVisit: "/disp/patient/visit",

    test: "/test"
}

const main = NewMiddleWare()
main.add(middlewareAuthRequire)

const disp = NewMiddleWare()
disp.add(middlewareAuthRequire)
const accessDisp = [
    {podr: 1, prava: 5}
]
disp.add(middlewareAccessRequire, {access: accessDisp})


const routes = [
    {
        path: linkDict.signin,
        element: <SignInPageState/>
    },
    {
        path: linkDict.start,
        element: main.middleware(<StartPage/>)
    },
    {
        path: linkDict.profile,
        element: disp.middleware(<ProfilePage/>)
    },
    {
        path: linkDict.findPatient,
        element: disp.middleware(<FoundPatient/>)
    },
    {
        path: linkDict.dispVisit,
        element: disp.middleware(<VisitPage/>)
    },
    {
        path: linkDict.patient,
        element: disp.middleware(<GetPatient/>)
    },


    {
        path: linkDict.test,
        element: disp.middleware(<TestPage/>)
    },
    {
        path: "*",
        element: disp.middleware(<StartPage/>)
    }
]
export default routes