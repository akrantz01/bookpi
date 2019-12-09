import React, { Component } from 'react';
import { HashRouter, Switch, Route } from 'react-router-dom';
import { ToastContainer, toast, Flip } from "react-toastify";

import { Users, Authentication, Shares, Files, Messages, Chats } from './api';

import Header from "./components/header";
import Footer from "./components/footer";

import Home from './pages/home';
import SignIn from './pages/sign-in';
import SignUp from "./pages/sign-up";
import NotFound from "./pages/not-found";
import Chat from "./pages/chat";

class Router extends Component {
    constructor(props) {
        super(props);

        this.state = {
            loggedIn: false,
            loading: true,
            user: {}
        };

        this.api = { Users, Authentication, Shares, Files, Messages, Chats };
    }

    componentDidMount() {
        Users.readSelf().then(data => {
            if (data.status === 200) this.setState({ loggedIn: true, user: data.data, loading: false });
            else if (data.status === 500) {
                toast.error(data.reason);
            }
            else this.setState({ loading: false });
        });
    }

    login() {
        this.setState({ loggedIn: true, loading: true });
        Users.readSelf().then(data => {
            if (data.status === 200) this.setState({ user: data.data, loading: false });
            else toast.error(data.reason);
        });
    }
    logout = () => this.setState({ loggedIn: false, user: {} });

    toggleLoading = () => this.setState({ loading: !this.state.loading });

    render() {
        if (this.state.loading) return (
            <div className="spinner-border text-info" style={{ position: "absolute", top: "40%", right: "50%" }} role="status">
                <span className="sr-only">Loading...</span>
            </div>
        );

        return (
            <HashRouter>
                <ToastContainer position="bottom-right" autoClose={4000} closeOnClick pauseOnHover draggable transition={Flip}/>
                <Header loggedIn={this.state.loggedIn} logout={this.logout.bind(this)}/>
                <main role="main" className="flex-shrink-0" style={{ marginTop: "40px" }}>
                    <Switch>
                        <Route path="/" exact><Home/></Route>
                        <Route path="/sign-in" exact><SignIn login={this.login.bind(this)} loggedIn={this.state.loggedIn} /></Route>
                        <Route path="/sign-up" exact><SignUp loggedIn={this.state.loggedIn} /></Route>
                        <Route path="/chat" exact><Chat loggedIn={this.state.loggedIn} username={this.state.user.username} /></Route>
                        <Route path="*"><NotFound/></Route>
                    </Switch>
                </main>
                <Footer/>
            </HashRouter>
        );
    }
}

export default Router;
