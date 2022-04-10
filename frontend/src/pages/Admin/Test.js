import React, {useState} from "react";

import * as appActions from "../../store/actions/application"

import {SideLoading} from "../../components/Loading";
import Notify, {notifyType} from "../../components/Notify";
import Modal, {BTN_CANCEL, BTN_CLOSE, BTN_NO, BTN_OK, BTN_SAVE, BTN_YES} from "../../components/Modal";


const Test = ({dispatch}) => {
  const [isLoading, setIsLoading] = useState(false)
  const [isOpenModal, setOpenModal] = useState(false)
  return <div>
    <div>
      <button onClick={_ => Notify(notifyType.INFO, "info")()}>info</button>
      <button onClick={_ => Notify(notifyType.SUCCESS, "success")()}>success</button>
      <button onClick={_ => Notify(notifyType.ERROR, "error")()}>error</button>
      <button onClick={_ => Notify(notifyType.WARNING, "warning")()}>warning</button>
    </div>
    <div>
      <button onClick={_ => setIsLoading(true)}>on</button>
      <button onClick={_ => setIsLoading(false)}>off</button>
      <SideLoading isLoading={isLoading}/>
    </div>
    <div>
      <button onClick={_ => dispatch(appActions.enableLoading())}>on</button>
      <button onClick={_ => dispatch(appActions.disableLoading())}>off</button>
    </div>
    <div>
      <button onClick={_ => dispatch(appActions.loadingAdd("test1"))}>on</button>
      <button onClick={_ => dispatch(appActions.loadingRemove("test1"))}>off</button>
    </div>
    <div>
      <button onClick={_ => setOpenModal(!isOpenModal)}>modal toggle</button>
      <Modal
        isOpen={isOpenModal}
        title={"Modal window"}
        body={<div>
            test modal
        </div>}
        btnNum={BTN_NO+BTN_OK+BTN_CANCEL+BTN_YES+BTN_CLOSE+BTN_SAVE}
        onNo={_=>alert("no")}
        onClose={_=>setOpenModal(false)}
      />
    </div>
  </div>
}

export default Test