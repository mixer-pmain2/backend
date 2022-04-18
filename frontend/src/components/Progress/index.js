import React, {memo} from "react";


const Progress = ({style}) => {

    return <div className="progress" style={{height: 5, ...style}}>
        <div className="progress-bar progress-bar-striped progress-bar-animated bg-danger"  aria-valuenow="100"
             aria-valuemin="0" aria-valuemax="100" style={{width: "100%"}}/>
    </div>
}

export default memo(Progress)