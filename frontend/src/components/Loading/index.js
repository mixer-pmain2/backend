import React from "react";

export const typeLoading = {
    PRIMARY: "primary",
    SECONDARY: "secondary",
    SUCCESS: "success",
    DANGER: "danger",
    WARNING: "warning",
    INFO: "info",
    LIGHT: "light",
    DARK: "dark"
}

const Loading = ({isLoading = false, type = typeLoading.PRIMARY}) => {
    if (!isLoading) return null

    return <div className={`spinner-border text-${type}`} role="status">
        <span className="visually-hidden">Loading...</span>
    </div>
}

export const SideLoading = (props) => {

    return <div style={{position: "absolute", right: 65, bottom: 65}}>
        <Loading {...props}/>
    </div>
}

export default Loading