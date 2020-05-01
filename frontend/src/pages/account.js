import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Users } from '../api';

class Account extends Component {
    constructor(props) {
        super(props);

        this.state = {
            name: "",
            password: "",
            confirmPassword: ""
        }
    }

    componentDidMount() {
        if (!this.props.loggedIn) this.props.history.push("/sign-in");
    }

    onNameInput = event => this.setState({ name: event.target.value });
    onPasswordInput = event => this.setState({ password: event.target.value });
    onConfirmPasswordInput = event => this.setState({ confirmPassword: event.target.value });

    onSubmit = () => {
        if (this.state.password !== "" || this.state.confirmPassword !== "" ) {
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
        }

        Users.update(this.state.name, this.state.password).then(data => {
            if (data.status === 200) {
                toast.success("Successfully updated your information");
                if (this.state.name !== "") this.props.updateName(this.state.name);
                this.setState({ name: "", password: "", confirmPassword: "" });
            }
            else if (data.status !== 401) toast.error(`Failed to update user: (${data.status}) ${data.reason}`);
        })
    };

    render() {
        return (
            <div className="container">
                <div className="card">
                    <div className="card-header">Welcome, {this.props.user.name}!</div>
                    <div className="card-body">
                        <h3 className="card-title">Manage your account</h3>
                        <h6 className="card-subtitle mb-2 text-muted">Update your name or password</h6>
                        <br/>

                        <form>
                            <div className="form-group">
                                <label htmlFor="updateName">Name</label>
                                <input type="text" className="form-control" id="updateName" placeholder={this.props.user.name} onChange={this.onNameInput.bind(this)} value={this.state.name}/>
                            </div>

                            <div className="form-row">
                                <div className="form-group col-md-6">
                                    <label htmlFor="updatePassword">Password:</label>
                                    <input type="password" className="form-control" id="updatePassword" onChange={this.onPasswordInput.bind(this)} value={this.state.password}/>
                                </div>

                                <div className="form-group col-md-6">
                                    <label htmlFor="updateConfirmPassword">Confirm Password:</label>
                                    <input type="text" className="form-control" id="updateConfirmPassword" onChange={this.onConfirmPasswordInput.bind(this)} value={this.state.confirmPassword}/>
                                </div>
                            </div>

                            <button type="submit" className="btn btn-primary" onClick={this.onSubmit.bind(this)}>Save</button>
                        </form>
                    </div>
                </div>
            </div>
        );
    }
}

Account.propTypes = {
    loggedIn: PropTypes.bool,
    history: PropTypes.object,
    user: PropTypes.shape({
        name: PropTypes.string,
        username: PropTypes.string
    }),
    updateName: PropTypes.func,
};

export default withRouter(Account);
