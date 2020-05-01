import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Link, withRouter } from 'react-router-dom';
import { Authentication } from "../api";
import {toast} from "react-toastify";

class Header extends Component {
    signOut() {
        Authentication.logout().then(data => {
            if (data.status === 200) this.props.logout();
            else if (data.status !== 401) toast.error(`Failed to logout: (${data.status}) ${data.reason}`);
        })
    }

    render() {
        let { location, loggedIn } = this.props;

        if (location.pathname === "/sign-in" || location.pathname === "/sign-up") return (<></>);

        return (
            <header>
                <nav className="navbar navbar-expand-lg navbar-dark fixed-top bg-dark">
                    <Link to="/" className="navbar-brand">BookPi</Link>
                    <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"/>
                    </button>

                    <div className="collapse navbar-collapse" id="navbarSupportedContent">
                        <ul className="navbar-nav mr-auto">
                            <li className="nav-item">
                                <Link className="nav-link" to="/">{ (location.pathname === "/") ? <b>Home</b> : "Home" }</Link>
                            </li>
                            { loggedIn && (
                                <>
                                    <li className="nav-item">
                                        <Link className="nav-link" to="/chat">{ (location.pathname === "/chat") ? <b>Chat</b> : "Chat" }</Link>
                                    </li>
                                    <li className="nav-item">
                                        <Link className="nav-link" to="/my-files">{ (location.pathname === "/my-files") ? <b>My Files</b> : "My Files" }</Link>
                                    </li>
                                    <li className="nav-item">
                                        <Link className="nav-link" to="/shared-with-me">{ (location.pathname === "/shared-with-me") ? <b>Shared with Me</b> : "Shared with Me" }</Link>
                                    </li>
                                    <li className="nav-item">
                                        <Link className="nav-link" to="/account">{ (location.pathname === "/account") ? <b>Account</b> : "Account" }</Link>
                                    </li>
                                </>
                            )}
                        </ul>

                        <form className="form-inline mr-sm-2">
                            { !loggedIn && (
                                <>
                                    <button className="btn btn-outline-success" type="button" onClick={() => this.props.history.push("/sign-in")}>Sign In</button>
                                    <span className="text-muted" style={{ marginLeft: "5px", marginRight: "5px" }}>|</span>
                                    <button className="btn btn-outline-primary" type="button" onClick={() => this.props.history.push("/sign-up")}>Sign Up</button>
                                </>
                            )}
                            { loggedIn && <button className="btn btn-outline-light" type="button" onClick={this.signOut.bind(this)}>Sign out</button> }
                        </form>
                    </div>
                </nav>
            </header>
        )
    }
}

Header.propTypes = {
    loggedIn: PropTypes.bool,
    logout: PropTypes.func,
    location: PropTypes.object
};

export default withRouter(Header);
