import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Link, withRouter } from 'react-router-dom';

class Header extends Component {
    render() {
        let { location, loggedIn } = this.props;
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
                            <li className="nav-item">
                                <Link className="nav-link" to="/files">{ (location.pathname === "/files") ? <b>Files</b> : "Files" }</Link>
                            </li>
                            <li className="nav-item">
                                <Link className="nav-link" to="/chat">{ (location.pathname === "/chat") ? <b>Chat</b> : "Chat" }</Link>
                            </li>
                        </ul>

                        <form className="form-inline mr-sm-2">
                            { !loggedIn && (
                                <>
                                    <button className="btn btn-outline-success" type="button" onClick={this.props.toggleLogin}>Sign In</button>
                                    <span className="text-muted" style={{ marginLeft: "5px", marginRight: "5px" }}>|</span>
                                    <button className="btn btn-outline-primary" type="button">Sign Up</button>
                                </>
                            )}
                            { loggedIn && <button className="btn btn-outline-light" type="button" onClick={this.props.toggleLogin}>Sign out</button> }
                        </form>
                    </div>
                </nav>
            </header>
        )
    }
}

Header.propTypes = {
    loggedIn: PropTypes.bool,
    toggleLogin: PropTypes.func,
    location: PropTypes.object
};

export default withRouter(Header);
