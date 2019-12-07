import React, { Component } from 'react';
import {HashRouter, Switch, Route} from 'react-router-dom';

import Header from "./components/header";
import Footer from "./components/footer";
import Home from './pages/home';

class Router extends Component {
    constructor(props) {
        super(props);

        this.state = {
            loggedIn: false,
            user: {}
        }
    }

    toggleLogin = () => this.setState({ loggedIn: !this.state.loggedIn });

    render() {
        return (
            <HashRouter>
                <Header loggedIn={this.state.loggedIn} toggleLogin={this.toggleLogin.bind(this)}/>
                <main role="main" className="flex-shrink-0" style={{ marginTop: "65px" }}>
                    <Switch>
                        <Route path="/" exact><Home/></Route>
                    </Switch>
                </main>
                <Footer/>
            </HashRouter>
        );
    }
}

export default Router;
