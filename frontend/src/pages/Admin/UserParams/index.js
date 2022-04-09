import React, {useState} from "react";
import UserAccess from "./UserAccess";
import UserFind from "./UserFind";

const UserParams = () => {
    const [selectUser, setSelectUser] = useState({})

    return <div>
        <UserFind selectUser={selectUser} setSelectUser={setSelectUser}/>
        {selectUser?.id && <UserAccess selectUser={selectUser}/>}
    </div>
}

export default UserParams