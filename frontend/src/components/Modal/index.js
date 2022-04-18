import React from "react";

export const BTN_CLOSE = 1;
export const BTN_SAVE = 2;
export const BTN_OK = 4;
export const BTN_CANCEL = 8;
export const BTN_YES = 16;
export const BTN_NO = 32;


const Modal = ({onClose, onSave, onOk, onNo, onYes, onCancel, isOpen = false, title = "", body, btnNum = 1, style={}}) => {

  if (!isOpen) {
    return null
  }

  return <div className="modal" tabIndex="-1" style={{display: "block", backgroundColor: "rgba(0,0,0,.2)"}}>
    <div className="modal-dialog" style={{opacity: 1, ...style}}>
      <div className="modal-content">
        <div className="modal-header">
          <h5 className="modal-title">{title}</h5>
          <button
            type="button"
            className="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"
            onClick={onClose}
          />
        </div>
        <div className="modal-body">
          {body}
        </div>
        <div className="modal-footer">
          {(btnNum & BTN_SAVE) > 0    && <button type="button" className="btn btn-primary" onClick={onSave}>Сохранить</button>}
          {(btnNum & BTN_OK) > 0      && <button type="button" className="btn btn-primary" onClick={onOk}>Ok</button>}
          {(btnNum & BTN_YES) > 0     && <button type="button" className="btn btn-primary" onClick={onYes}>Да</button>}
          {(btnNum & BTN_NO) > 0      && <button type="button" className="btn btn-outline-primary" onClick={onNo}>Нет</button>}
          {(btnNum & BTN_CANCEL) > 0  && <button type="button" className="btn btn-outline-danger" onClick={onCancel}>Отмена</button>}
          {(btnNum & BTN_CLOSE) > 0   && <button type="button" className="btn btn-secondary" onClick={onClose}>
            Закрыть
          </button>}
        </div>
      </div>
    </div>
  </div>
}

export default Modal