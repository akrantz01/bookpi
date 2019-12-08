import React, { Component } from 'react';

class NotFound extends Component {
    render() {
        return (
            <div className="container">
                <div className="jumbotron">
                    <h1 className="display-4">404 <span className="text-muted">|</span> Not Found</h1>
                    <hr className="my-4"/>
                    <p className="lead">The page you requested does not exist. Please check the URL and try again.</p>
                </div>
            </div>
        )
    }
}

export default NotFound;
