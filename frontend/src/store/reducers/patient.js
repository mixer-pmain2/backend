import {createSlice} from "@reduxjs/toolkit";
import {saveToStore} from "./_save";

const storeName = "patient"

const initialState = localStorage.getItem(storeName) ? JSON.parse(localStorage.getItem(storeName)) : {}

export const patientStore = createSlice({
    name: storeName,
    initialState,
    reducers: {
        select: (state, action) => {
            state = {
                ...action.payload
            }
            saveToStore(state, storeName)
            return state
        },
        reset: state => {
            state = {}
            saveToStore(state, storeName)
            return state
        },
        setUchet: (state, action) => {
            state = {
                ...state,
                uchet: action.payload
            }
            saveToStore(state, storeName)
            return state
        },
        setVisits: (state, action) => {
            state = {
                ...state,
                visit: action.payload
            }
            saveToStore(state, storeName)
            return state
        },
        setHospital: (state, action) => {
            state = {
                ...state,
                hospital: action.payload
            }
            saveToStore(state, storeName)
            return state
        }
    }
})

export const {select, reset, setUchet, setVisits, setHospital} = patientStore.actions
export default patientStore.reducer