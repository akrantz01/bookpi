import React, { Component } from 'react';
import {HashRouter, Switch, Route} from 'react-router-dom';
import { ToastContainer, Flip } from "react-toastify";

import Header from "./components/header";
import Footer from "./components/footer";

import Home from './pages/home';
import SignIn from './pages/sign-in';

class Router extends Component {
    constructor(props) {
        super(props);

        this.state = {
            loggedIn: false,
            user: {},
        }
    }

    toggleLogin = () => this.setState({ loggedIn: !this.state.loggedIn });

    render() {
        return (
            <HashRouter>
                <ToastContainer position="top-right" autoClose={4000} newestOnTop closeOnClick pauseOnHover draggable transition={Flip}/>
                <Header loggedIn={this.state.loggedIn} toggleLogin={this.toggleLogin.bind(this)}/>
                <main role="main" className="flex-shrink-0" style={{ marginTop: "40px" }}>
                    <Switch>
                        <Route path="/" exact><Home/></Route>
                        <Route path="/sign-in" exact><SignIn/></Route>
                    </Switch>
                </main>
                <Footer/>
            </HashRouter>
        );
    }
}

export default Router;
