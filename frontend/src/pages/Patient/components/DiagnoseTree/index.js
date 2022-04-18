import React, {useEffect, useState} from "react";
import {connect} from "react-redux";

import * as appActions from "store/actions/application"
import * as apiSpr from "api/spr"
import {TreeView} from "components/Tree";
import Modal, {BTN_NO, BTN_YES} from "components/Modal";

export const getTypeDiagsModal = (uch) => {
  switch (uch) {
    case 11:
    case 12:
      return 1
    case 13:
    case 14:
    case 15:
      return 0
    case 16:
      return 2
    default:
      return 0
  }
}

const DiagnoseTree = ({dispatch, onSelect = console.log, type = 1}) => {
  const [diags, setDiags] = useState([])
  const [selected, setSelected] = useState("")
  const [isOpenModalDiag, setIsOpenModalDiag] = useState(false)
  const [isOpenConfirm, setIsOpenConfirm] = useState(false)

  const mapper = (data) => data.map((v, i) => ({
    title: v.title,
    value: v.diag,
    isNode: v.haveChildren
  }))

  const diagFilter = (diags, isMain = true) => {
    diags = diags || []
    /*
    0 - все
    1 - психиатрия + z
    2 - наркология + z
    3 - соматика
    4 - психиатрия + наркология + z
    5 - психиатрия
     */
    if (isMain) {
      if (type === 0) return diags
      if (type === 1) return diags.filter(v => ["F00-F99", "Z00-Z99"].indexOf(v.diag) + 1)
      if (type === 2) return diags.filter(v => ["F00-F99", "Z00-Z99"].indexOf(v.diag) + 1)
      if (type === 3) return diags.filter(v => !(["F00-F99"].indexOf(v.diag) + 1))
      if (type === 4) return diags.filter(v => ["F00-F99", "Z00-Z99"].indexOf(v.diag) + 1)
      if (type === 5) return diags.filter(v => ["F00-F99"].indexOf(v.diag) + 1)
    }
    if (type === 1) return diags.filter(v => (v.diag.length <= 6) ? true : !(["F10-F19"].indexOf(v.diag) + 1))
    if (type === 2) return diags.filter(v => (v.diag.length <= 6) ? true : ((["F10-F19"].indexOf(v.diag) + 1) || v.head === "Z00-Z99"))
    if (type === 5) return diags.filter(v => (v.diag.length <= 6) ? true : !(["F10-F19"].indexOf(v.diag) + 1))
    return diags
  }

  const getDiagMain = (diag = "", mapper, setData) => {
    const nameLoaded = "spr_diag_main_" + diag
    dispatch(appActions.loadingAdd(nameLoaded))
    apiSpr.getSprDiag({diag})
      .then(res => {
        setData(mapper(diagFilter(res)))
      })
      .finally(() => dispatch(appActions.loadingRemove(nameLoaded)))
  }

  const getDiag = (diag = "", mapper, setData) => {
    const nameLoaded = "spr_diag_" + diag
    dispatch(appActions.loadingAdd(nameLoaded))
    apiSpr.getSprDiag({diag})
      .then(res => {
        setData(mapper(diagFilter(res, false)))
      })
      .finally(() => dispatch(appActions.loadingRemove(nameLoaded)))
  }

  const handleSelect = (diag) => {
    setSelected(diag)
    setIsOpenConfirm(true)
  }

  useEffect(() => {
    getDiagMain("", mapper, setDiags)
  }, [])

  return <div>
    <TreeView
      data={diags}
      getData={(v, setData) => getDiag(v, mapper, setData)}
      onSelect={handleSelect}
    />

    <Modal
      title={"Вы выбрали диагноз " + selected}
      body={<div>
        <span>{selected} Это правильно?</span>
      </div>}
      btnNum={BTN_YES + BTN_NO}
      onYes={() => onSelect(selected)}
      onNo={() => setIsOpenConfirm(false)}
      isOpen={isOpenConfirm}
      onClose={() => setIsOpenConfirm(false)}
    />
  </div>
}

export default connect(state => ({
  application: state.application
}))(DiagnoseTree)