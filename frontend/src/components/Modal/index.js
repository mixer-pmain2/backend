import React from "react";

const modalBtn = {
    CLOSE: 0,
    SAVE: 1
}

const Modal = ({onClose, isOpen = false, title = "", body}) => {

    if (!isOpen) {
        return null
    }
    console.log(isOpen)
    return <div className="modal" tabIndex="-1" style={{display: "block"}}>
        <div className="modal-dialog">
            <div className="modal-content">
                <div className="modal-header">
                    <h5 className="modal-title">Modal title</h5>
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
                    <button type="button" className="btn btn-secondary" data-bs-dismiss="modal" onClick={onClose}>
                        Закрыть
                    </button>
                    <button type="button" className="btn btn-primary">Save changes</button>
                </div>
            </div>
        </div>
    </div>
}

export default Modal