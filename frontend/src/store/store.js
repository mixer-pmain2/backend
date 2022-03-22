import { configureStore } from '@reduxjs/toolkit'

import userReducer from "./reducers/user";
import patientReducer from "./reducers/patient"
import appReducer from "./reducers/application"

export const store = configureStore({
    reducer: {
        user: userReducer,
        patient: patientReducer,
        application: appReducer,
    },
})