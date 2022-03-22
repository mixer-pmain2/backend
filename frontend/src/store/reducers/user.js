import {createSlice} from "@reduxjs/toolkit";
import {saveToStore} from "./_save";

const storeName = "user"

const initialState = localStorage.getItem(storeName) ? JSON.parse(localStorage.getItem(storeName)) : {
    isAuth: false
}

export const userStore = createSlice({
    name: storeName,
    initialState,
    reducers: {
        login: (state, action) => {
            state = {
                ...state,
                isAuth: true,
                ...action.payload
            }
            saveToStore(state, storeName)
            return state
        },
        logout: state => {
            state = {
                isAuth: false
            }
            saveToStore(state, storeName)
            return state
        },
        setToken: (state, action) => {
            state = {
                ...state,
                token: action.payload.token
            }
            saveToStore(state, storeName)
            return state
        },
        setPrava: (state, action) => {
            state = {
                ...state,
                prava: action.payload
            }
            saveToStore(state, storeName)
            return state
        },
        setCurrentPodr: (state, action) => {
            state = {
                ...state,
                podr: action.payload
            }
            saveToStore(state, storeName)
            return state
        }
    }
})

export const { login, logout, setToken, setPrava, setCurrentPodr } = userStore.actions
export default userStore.reducer