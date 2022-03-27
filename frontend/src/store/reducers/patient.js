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
        }
    }
})

export const {select, reset, setUchet} = patientStore.actions
export default patientStore.reducer