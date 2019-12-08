import React from 'react';
import ReactDOM from 'react-dom';
import * as serviceWorker from './serviceWorker';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'react-toastify/dist/ReactToastify.min.css';

import Router from './router';

ReactDOM.render(<Router/>, document.getElementById('root'));

serviceWorker.unregister();
