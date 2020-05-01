import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFileDownload, faFileAlt } from "@fortawesome/free-solid-svg-icons";
import { toast } from 'react-toastify';
import { Shares } from '../api';

class Shared extends Component {
    constructor(props) {
        super(props);

        this.state = {
            loading: false,
            shares: []
        }
    }

    componentDidMount() {
        this.refresh();
    }

    refresh() {
        this.setState({ loading: true });
        Shares.list()
            .then(data => {
                if (data.status !== 200) toast.error(`Failed to load shared files: (${data.status}) ${data.reason}`);
                else {
                    let shares = {};
                    for (let share of data.data) {
                        let [author, path] = share.split('/', 2)
                        if (!shares.hasOwnProperty(author)) shares[author] = [path];
                        else shares[author].push(path);
                    }
                    this.setState({ shares });
                }
            })
            .finally(() => this.setState({ loading: false }));
    }

    download = (user, file) => () => Shares.download(user, file)
        .then(data => {
            if (data.status !== 200) toast.error(`Failed to download file: (${data.status}) ${data.reason}`);
        });

    render() {
        if (!this.props.loggedIn) this.props.history.push("/sign-in");

        return (
            <div className="container">
                <div className="card" style={{ height: "85vh" }}>
                    <div className="card-header">
                        <div className="row">
                            <div className="col-sm">
                                <h3 className="card-title">Shared with Me</h3>
                            </div>
                            <div className="col-sm text-right">
                                <button type="button" className="btn btn-outline-primary" onClick={this.refresh.bind(this)}>Refresh</button>
                            </div>
                        </div>
                    </div>
                    <div className="card-body" style={{ overflow: "auto", paddingTop: "0px" }}>
                        { this.state.loading && (
                            <div className="spinner-border" role="status">
                                <span className="sr-only">Loading...</span>
                            </div>
                        )}
                        { !this.state.loading && Object.keys(this.state.shares).length === 0 && <h6>No files were shared with you!</h6>}
                        { !this.state.loading && Object.keys(this.state.shares).length !== 0 && Object.keys(this.state.shares).map(author => (
                            <>
                                <div className="card" style={{ marginTop: "0.25rem", marginBottom: "0.25rem" }}>
                                    <div className="card-body" style={{ padding: "0.625rem", marginTop: "0.375rem" }}>
                                        <h5>From: <code>{author}</code></h5>
                                    </div>
                                </div>
                                { this.state.shares[author].map(file => (
                                    <div className="card" style={{ marginTop: "0.25rem", marginBottom: "0.25rem" }}>
                                        <div className="card-body" style={{ padding: "0.75rem" }}>
                                            <div className="container">
                                                <div className="row">
                                                    <div className="col-sm" style={{ paddingTop: "0.125rem" }}>
                                                        <FontAwesomeIcon icon={faFileAlt} style={{ fontSize: "0.75rem" }}/> &nbsp;{file}
                                                    </div>
                                                    <div className="col-sm text-right">
                                                        <button type="button" className="btn btn-outline-success btn-sm" style={{ fontSize: "0.75rem", marginRight: "0.25rem" }}
                                                                title="Download" onClick={this.download(author, file)}><FontAwesomeIcon icon={faFileDownload}/></button>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </>
                        ))}
                    </div>
                </div>
            </div>
        )
    }
}

Shared.propTypes = {
    loggedIn: PropTypes.bool,
    history: PropTypes.object,
};

export default withRouter(Shared);
