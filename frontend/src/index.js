import React from 'react';
import ReactDOM from 'react-dom';
import {HashRouter as Router, Switch, Route} from 'react-router-dom';
import * as serviceWorker from './serviceWorker';

import Header from "./components/header";
import Footer from "./components/footer";
import Home from './pages/home';

import 'bootstrap/dist/css/bootstrap.min.css';

ReactDOM.render(<Router>
    <Header/>
    <main role="main" className="flex-shrink-0" style={{ marginTop: "65px" }}>
        <Switch>
            <Route path="/" exact><Home/></Route>
        </Switch>
    </main>
    <Footer/>
</Router>, document.getElementById('root'));

serviceWorker.unregister();
