import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import reportWebVitals from './reportWebVitals';

ReactDOM.render(
    <App/>,
    document.getElementById('root')
);

reportWebVitals();

require("babel-core/register");
require("babel-polyfill");