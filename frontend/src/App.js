import React from "react";
import {Provider} from 'react-redux'
import {BrowserRouter as Router, Route, Routes} from "react-router-dom";


import {store} from "./store/store";

import routes from "./routes";
import Initialisation from "./initialisation";

function App() {

    return (
        <Provider store={store}>
            <Router>
                <Routes>
                    {routes.map((v, i) =>
                        <Route key={i} path={v.path} element={v.element}/>)
                    }
                </Routes>
            </Router>

            <Initialisation/>
        </Provider>
    );
}

export default App;
