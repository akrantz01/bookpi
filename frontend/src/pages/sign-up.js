import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import { Link, withRouter } from 'react-router-dom';
import { Authentication } from "../api";

import '../style/user-form.css'

class SignUp extends Component {
    constructor(props) {
        super(props);

        this.state = {
            name: "",
            username: "",
            password: "",
            confirmPassword: ""
        }
    }

    componentDidMount() {
        if (this.props.loggedIn) this.props.history.push("/");
    }

    onNameInput = event => this.setState({ name: event.target.value });
    onUsernameInput = event => this.setState({ username: event.target.value });
    onPasswordInput = event => this.setState({ password: event.target.value });
    onConfirmPasswordInput = event => this.setState({ confirmPassword: event.target.value });

    onSubmit = () => {
        if (this.state.password !== this.state.confirmPassword) return toast.error("Passwords do not match");
        if (!this.state.password.match(/^(?=.{8,})(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!-/:-@[-_]).*$/)) return toast.error(<div>
            <p>Your password must meet the following requirements:</p>
            <ul>
                <li>at least 8 characters long</li>
                <li>1 lowercase character</li>
                <li>1 uppercase character</li>
                <li>1 special character</li>
                <li>1 number</li>
            </ul>
        </div>);

        Authentication.register(this.state.username, this.state.password, this.state.name).then(data => {
            if (data.status === 200) {
                toast.success("Successfully signed up, you may now login");
                this.props.history.push("/sign-in");
            }
            else toast.error(data.reason);
        });
    };

    render() {
        return (
            <div style={{ alignItems: "center" }} className="text-center">
                <form className="form-signin">
                    <h1 className="h3 mb-3 font-weight-normal">Sign up for BookPi</h1>

                    <label htmlFor="name" className="sr-only">Name</label>
                    <input type="text" id="name" className="form-control form-top" placeholder="Name" onInput={this.onNameInput.bind(this)} required autoFocus/>

                    <label htmlFor="username" className="sr-only">Username</label>
                    <input type="text" id="username" className="form-control form-bottom" placeholder="Username" onInput={this.onUsernameInput.bind(this)} required/>

                    <label htmlFor="password" className="sr-only">Password</label>
                    <input type="password" id="password" className="form-control form-top" placeholder="Password" onInput={this.onPasswordInput.bind(this)} required/>

                    <label htmlFor="password-confirm" className="sr-only">Confirm Password</label>
                    <input type="password" id="password-confirm" className="form-control form-bottom" placeholder="Confirm your password" onInput={this.onConfirmPasswordInput.bind(this)} required/>

                    <button className="btn btn-lg btn-primary btn-block" type="submit" onClick={this.onSubmit.bind(this)}>Sign Up</button>
                    <br/>
                    <p className="text-muted">Already have an account? <Link to="/sign-in">Sign in!</Link></p>
                    <p className="text-muted"><Link to="/" style={{ color: "#6c757d", fontSize: "14px" }}>Home</Link></p>
                </form>
            </div>
        );
    }
}

SignUp.propTypes = {
    history: PropTypes.object,
    loggedIn: PropTypes.bool
};

export default withRouter(SignUp);
