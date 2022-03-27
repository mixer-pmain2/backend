import React from "react";
import StartPage from "./pages/Start";
import VisitPage from "./pages/Patient/Dispanser/modules/Visit/index";


import GetPatient from "./pages/Patient";
import ProfilePage from "./pages/Profile";
import middlewareAuthRequire from "./middleware/authRequier";
import NewMiddleWare from "./middleware";
import middlewareAccessRequire from "./middleware/accessRequier";
import SignInPageState from "./pages/Signin";
import AdminPage from "./pages/Admin";

import * as access from "./configs/access";
import Page404 from "./pages/ErrorPage/404";
import FindPatient from "./pages/Patient/Find";

export const linkDict = {
    start: "/",

    profile: "/profile",
    signin: "/signin",

    disp: "/disp",
    patient: "/patient/:id",
    findPatient: "/patient/find",
    dispVisit: "/disp/patient/visit",

    admin: "/admin",
}

const main = NewMiddleWare()
main.add(middlewareAuthRequire)

const disp = NewMiddleWare()
disp.add(middlewareAccessRequire, {access: access.accessDispPage})
disp.add(middlewareAuthRequire)

const admin = NewMiddleWare()
admin.add(middlewareAuthRequire)
admin.add(middlewareAccessRequire, {access: access.accessAdminPage})


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
        element: disp.middleware(<FindPatient/>)
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
        path: linkDict.admin,
        element: admin.middleware(<AdminPage/>)
    },
    {
        path: "*",
        element: <Page404/>
    }
]
export default routes