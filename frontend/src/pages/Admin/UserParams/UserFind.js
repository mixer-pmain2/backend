import React, {useState} from "react";

import * as userApi from "../../../api/user"

const UserFind = () => {
    const [selectUser, setSelectUser] = useState({})

    const onReset = () => {

    }

    const findUserById = (id) => {
        userApi.getUser({id})
            .then(res => {
                if (res.id > 0) {
                    setSelectUser(res)
                }
            })
    }

    const handleFindUser = (e) => {
        e.preventDefault()
        const [userId] = e.target
        if (Number(userId.value) > 0) {
            findUserById(userId.value)
        }
    }
    console.log(selectUser)
    return <div>
        <form onSubmit={handleFindUser} autoComplete="off">
            <div className="mb-3 d-flex flex-row">
                <a href="#" className="btn btn-outline-secondary" style={{marginRight: 10}} onClick={onReset}>Сброс</a>
                <div className="input-group w-25" style={{marginRight: 10}} >
                    <span className="input-group-text">Шифр</span>
                    <input
                        type="text"
                        className="form-control"
                        id="userId"
                        name="userId"
                        autoFocus={true}
                        maxLength={10}
                        onChange={e => {
                        }}
                        value={selectUser?.id}
                    />

                </div>
                <div className="input-group w-100" style={{marginRight: 10}}>
                    <span className="input-group-text">Ф.И.О.</span>
                    <input
                        type="text"
                        className="form-control"
                        id="fio"
                        name="fio"
                        maxLength={255}
                        onChange={e => {
                        }}
                        value={(selectUser?.lname??"")+" "+(selectUser?.fname??"")+" "+(selectUser?.sname??"")}
                    />
                </div>
                <button className="input-group-text btn btn-outline-primary">Найти</button>
            </div>
        </form>
    </div>
}

export default UserFind