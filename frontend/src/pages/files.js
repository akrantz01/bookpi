import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Files, Shares } from '../api';
import {DirectoryEntry, FileEntry} from "../components/entries";

const style = {
    breadcrumb: {
        paddingTop: "0rem",
        paddingBottom: "0rem",
        paddingLeft: "0.375rem",
        paddingRight: "0.375rem",
    }
};

class FileManager extends Component {
    constructor(props) {
        super(props);

        this.state = {
            loadingFiles: true,
            currentDirectory: "",
            children: []
        }
    }

    componentDidMount() {
        this.refreshFilesList();
    }

    upDirectory = idx => () => {
        this.setState({ currentDirectory: this.state.currentDirectory.split("/").slice(0, idx+1).join("/") });
        setTimeout(() => this.refreshFilesList(), 50);
    }

    downDirectory = newDir => () => {
        this.setState({ currentDirectory: `${this.state.currentDirectory}/${newDir}` });
        setTimeout(() => this.refreshFilesList(), 50);
    }

    refreshFilesList() {
        this.setState({ loadingFiles: true });
        Files.read(this.state.currentDirectory).then(files => this.setState({ loadingFiles: false, children: files.data.children }));
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
        if (!this.props.loggedIn) this.props.history.push("/sign-in");

        return (
            <div className="container">
                <div className="row">
                    <div className="col-sm-4">
                        <div className="card" style={{ height: "85vh" }}>
                            <div className="card-body">
                                <h3 className="card-title">Files</h3>
                                <button type="button" className="btn btn-sm btn-outline-success">Upload</button>
                                <hr className="my-4"/>
                                <div className="list-group list-group-flush">
                                    <button type="button" className="list-group-item list-group-item-action">My Files</button>
                                    <button type="button" className="list-group-item list-group-item-action">Shared Files</button>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="col-sm-8">
                        <div className="card" style={{ height: "85vh" }}>
                            <nav aria-label="breadcrumb">
                                <ol className="breadcrumb">{ this.generateBreadcrumbs() }</ol>
                            </nav>
                            <div className="card-body" style={{ overflow: "auto", paddingTop: "0px" }}>
                                { this.state.loadingFiles && (
                                    <div className="spinner-border" role="status">
                                        <span className="sr-only">Loading...</span>
                                    </div>
                                )}
                                { !this.state.loadingFiles && !this.state.children && <h6>You have no files!</h6> }
                                { !this.state.loadingFiles && this.state.children && this.state.children.map(
                                    data => (data.directory) ?
                                        <DirectoryEntry key={Math.random().toString()} directory={data} onClick={this.downDirectory(data.name)}/> :
                                        <FileEntry key={Math.random().toString()} file={data}/>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

FileManager.propTypes = {
    loggedIn: PropTypes.bool,
    history: PropTypes.object,
    username: PropTypes.string,
};

export default withRouter(FileManager);
