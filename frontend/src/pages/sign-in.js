import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import { Link, withRouter } from 'react-router-dom';
import { Authentication } from "../api";

import '../style/user-form.css'

class SignIn extends Component {
    constructor(props) {
        super(props);

        this.state = {
            username: "",
            password: ""
        }
    }

    componentDidMount() {
        if (this.props.loggedIn) this.props.history.push("/");
    }

    onUsernameInput = event => this.setState({ username: event.target.value });
    onPasswordInput = event => this.setState({ password: event.target.value });

    onSubmit = () => Authentication.login(this.state.username, this.state.password).then(data => {
        if (data.status === 200) {
            this.props.login();
            toast.success("Successfully logged in");
            this.props.history.push("/");
        }
        else toast.error(data.reason);
    });

    render() {
        return (
            <div style={{ alignItems: "center" }} className="text-center">
                <form className="form-signin">
                    <h1 className="h3 mb-3 font-weight-normal">Sign in to BookPi</h1>

                    <label htmlFor="username" className="sr-only">Username</label>
                    <input type="text" id="username" className="form-control form-top" placeholder="Username" onInput={this.onUsernameInput.bind(this)} required autoFocus/>

                    <label htmlFor="password" className="sr-only">Password</label>
                    <input type="password" id="password" className="form-control form-bottom" placeholder="Password" onInput={this.onPasswordInput.bind(this)} required/>

                    <button className="btn btn-lg btn-primary btn-block" type="submit" onClick={this.onSubmit.bind(this)}>Sign In</button>
                    <br/>
                    <p className="text-muted">Don't have an account? <Link to="/sign-up">Sign up!</Link></p>
                    <p className="text-muted"><Link to="/" style={{ color: "#6c757d", fontSize: "14px" }}>Home</Link></p>
                </form>
            </div>
        );
    }
}

SignIn.propTypes = {
    login: PropTypes.func,
    history: PropTypes.object,
    loggedIn: PropTypes.bool
};

export default withRouter(SignIn);
