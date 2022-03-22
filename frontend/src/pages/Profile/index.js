import React from "react";
import {connect} from "react-redux";

const ProfilePage = ({dispatch, user}) => {

    console.log(user)
    return <div>

    </div>
}

export default connect(state => ({
    user: state.user
}))(ProfilePage)