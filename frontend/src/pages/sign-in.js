import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import { withRouter } from 'react-router-dom';
import { Authentication } from "../api";

import '../style/sign-in.css'

class SignIn extends Component {
    constructor(props) {
        super(props);

        this.state = {
            username: "",
            password: ""
        }
    }

    onUsernameInput = event => this.setState({ username: event.target.value });
    onPasswordInput = event => this.setState({ password: event.target.value });

    onSubmit = () => Authentication.login(this.state.username, this.state.password).then(data => {
        if (data.status === 200) {
            this.props.toggleLogin();
            this.props.history.back();
            toast.success("Successfully logged in");
        }
        else toast.error(data.reason);
    });

    render() {
        return (
            <div className="text-center">
                <form className="form-signin">
                    <h1 className="h3 mb-3 font-weight-normal">Please sign in</h1>

                    <label htmlFor="username" className="sr-only">Username</label>
                    <input type="text" id="username" className="form-control form-top" placeholder="Username" onInput={this.onUsernameInput.bind(this)} required autoFocus/>

                    <label htmlFor="password" className="sr-only">Password</label>
                    <input type="password" id="password" className="form-control form-bottom" placeholder="Password" onInput={this.onPasswordInput.bind(this)} required/>

                    <button className="btn btn-lg btn-primary btn-block" type="submit" onClick={this.onSubmit.bind(this)}>Sign In</button>
                </form>
            </div>
        );
    }
}

SignIn.propTypes = {
    toggleLogin: PropTypes.func,
    history: PropTypes.object
};

export default withRouter(SignIn);
