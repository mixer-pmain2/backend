import React from "react";


const Test = ({children}) => {

    console.log("middleware Test")

    return <>{children}</>
}

const middlewareTest = (next, params) => {
    return <Test>
        {next}
    </Test>
}

export default middlewareTest