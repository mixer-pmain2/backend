import React from "react";
import {useNavigate} from "react-router-dom";
import {linkDict} from "../../routes";

const NotAccessedPage = () => {

    const navigate = useNavigate()

    return <div>
        <div>
            <span className="text-body">Доступ запрещен</span>
        </div>
        <div>
            <a href="#" onClick={_ => navigate(linkDict.start)}>На галвную</a>
        </div>
    </div>
}

export default NotAccessedPage