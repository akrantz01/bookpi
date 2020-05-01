import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Files as FilesApi } from '../api';
import Entry from "../components/entry";

const style = {
    breadcrumb: {
        paddingTop: "0rem",
        paddingBottom: "0rem",
        paddingLeft: "0.375rem",
        paddingRight: "0.375rem",
    }
};

class Files extends Component {
    constructor(props) {
        super(props);

        this.state = {
            loading: true,
            currentDirectory: "",
            children: []
        }
    }

    componentDidMount() {
        if (!this.props.loggedIn) return this.props.history.push("/sign-in");
        this.refresh();
    }

    upDirectory = idx => () => {
        this.setState({ currentDirectory: this.state.currentDirectory.split("/").slice(0, idx+1).join("/") });
        setTimeout(() => this.refresh(), 50);
    }

    downDirectory = newDir => () => {
        this.setState({ currentDirectory: `${this.state.currentDirectory}/${newDir}` });
        setTimeout(() => this.refresh(), 50);
    }

    refresh() {
        this.setState({ loading: true });
        FilesApi.read(this.state.currentDirectory)
            .then(data => {
                if (data.status !== 200 && data.status !== 401) toast.error(`Failed to load files: (${data.status}) ${data.reason}`);
                else this.setState({ children: data.data.children })
            })
            .finally(() => this.setState({ loading: false }));
    }

    generateBreadcrumbs() {
        if (this.state.currentDirectory === "") return <li className="breadcrumb-item active" aria-current="page">
            <button type="button" className="btn btn-link" style={style.breadcrumb} disabled={true}>Home</button></li>;

        let paths = this.state.currentDirectory.split("/");
        return paths.map((path, index) => <li key={index}
                                              className={"breadcrumb-item" + ((index === paths.length - 1) ? " active" : "")}
                                              aria-current={ (index === paths.length - 1) ? "page" : "false" }>
            <button type="button" className="btn btn-link" style={style.breadcrumb} disabled={index === paths.length - 1}
                    onClick={this.upDirectory(index)}>{ (index === 0) ? "Home" : path }</button>
        </li> );
    }

    render() {
        return (
            <div className="container">
                <div className="card" style={{ height: "85vh" }}>
                    <div className="card-header">
                        <div className="row" style={{ marginBottom: "1rem" }}>
                            <div className="col-sm">
                                <h3 className="card-title">Files</h3>
                            </div>
                            <div className="col-sm text-right">
                                <div className="btn-group" role="group" aria-label="File operations">
                                    <button type="button" className="btn btn-outline-primary" onClick={this.refresh.bind(this)}>Refresh</button>
                                    <button type="button" className="btn btn-outline-success">Upload</button>
                                    <button type="button" className="btn btn-outline-info">New Folder</button>
                                </div>
                            </div>
                        </div>

                        <nav aria-label="breadcrumb">
                            <ol className="breadcrumb">{ this.generateBreadcrumbs() }</ol>
                        </nav>
                    </div>
                    <div className="card-body" style={{ overflow: "auto", paddingTop: "0px" }}>
                        { this.state.loading && (
                            <div className="spinner-border" role="status">
                                <span className="sr-only">Loading...</span>
                            </div>
                        )}
                        { !this.state.loading && !this.state.children && <h6>You have no files!</h6> }
                        { !this.state.loading && this.state.children &&
                        this.state.children.map(data => <Entry data={data} onClick={this.downDirectory(data.name)} refresh={this.refresh.bind(this)}
                                                               currentDirectory={this.state.currentDirectory} key={data.name}/>)}
                    </div>
                </div>
            </div>
        )
    }
}

Files.propTypes = {
    loggedIn: PropTypes.bool,
    history: PropTypes.object,
};

export default withRouter(Files);
