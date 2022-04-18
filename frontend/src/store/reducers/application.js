import {createSlice} from "@reduxjs/toolkit";
import {saveToStore} from "./_save";

const storeName = "app"

const initialState = localStorage.getItem(storeName) ? JSON.parse(localStorage.getItem(storeName)) : {
  loading: false,
  loadingList: [],
  spr: {},
  params: null
}

export const applicationStore = createSlice({
  name: storeName,
  initialState,
  reducers: {
    loadingListAdd: (state, action) => {
      const l = state?.loadingList ? state?.loadingList : []
      state = {
        ...state,
        loadingList: [...l, action.payload]
      }
      return state
    },
    loadingListRemove: (state, action) => {
      state = {
        ...state,
        loadingList: state.loadingList.filter(v => v !== action.payload)
      }
      return state

    },
    loadingListReset: (state) => {
      state = {
        ...state,
        loadingList: []
      }
      return state
    },
    loadingEnable: (state) => {
      state = {
        ...state,
        loading: true
      }
      return state
    },
    loadingDisable: (state) => {
      state = {
        ...state,
        loading: false
      }
      return state
    },
    setSprPodr: (state, action) => {
      state = {
        ...state,
        spr: {
          ...state.spr,
          unit: action.payload
        }
      }
      saveToStore(state, storeName)
      return state
    },
    setSprPrava: (state, action) => {
      state = {
        ...state,
        spr: {
          ...state.spr,
          access: action.payload
        }
      }
      saveToStore(state, storeName)
      return state
    },
    setSprVisit: (state, action) => {
      state = {
        ...state,
        spr: {
          ...state.spr,
          visit: action.payload
        }
      }
      saveToStore(state, storeName)
      return state
    },
    setParams: (state, action) => {
      state = {
        ...state,
        params: action.payload
      }
      saveToStore(state, storeName)
      return state
    }
  }
})

export const {
  loadingEnable,
  loadingDisable,
  setSprPodr,
  setSprPrava,
  loadingListAdd,
  loadingListRemove,
  loadingListReset,
  setSprVisit,
  setParams
} = applicationStore.actions
export default applicationStore.reducer