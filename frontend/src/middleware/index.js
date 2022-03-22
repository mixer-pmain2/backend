const NewMiddleWare = () => {
    let listMiddleWare = []

    return {
        add: (mf, params = {}) => listMiddleWare.push({mw: mf, params}),
        middleware: (children) => {
            let nextHandler = children

            listMiddleWare.map(v => {
                return nextHandler = v.mw(nextHandler, v.params)
            })

            return nextHandler
        }
    }
}


export default NewMiddleWare
