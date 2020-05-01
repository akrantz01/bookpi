import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFileDownload, faFolder, faFileAlt, faEllipsisV, faShareAlt } from "@fortawesome/free-solid-svg-icons";
import ReactModal from 'react-modal';
import { toast } from 'react-toastify';
import { Files, Shares } from '../api';

export default class Entry extends Component {
    constructor(props) {
        super(props);

        this.state = {
            newName: "",
            newPath: "",
            shareTo: "",
            renameModalOpen: false,
            moveModalOpen: false,
            shareModalOpen: false,
            loading: false,
        }
    }

    toggleRenameModal = () => this.setState({ renameModalOpen: !this.state.renameModalOpen });
    toggleMoveModal = () => this.setState({ moveModalOpen: !this.state.moveModalOpen });
    toggleShareModal = () => this.setState({ shareModalOpen: !this.state.shareModalOpen });

    onNameChange = e => this.setState({ newName: e.target.value });
    onPathChange = e => this.setState({ newPath: e.target.value });
    onUserChange = e => this.setState({ shareTo: e.target.value });

    download = () => Files.read(`${this.props.currentDirectory}/${this.props.data.name}`, true)
        .then(data => {
            if (data.status !== 200 && data.status !== 401) toast.error(`Failed to download file: (${data.status}) ${data.reason}`);
        });
    remove = () => Files.delete(`${this.props.currentDirectory}/${this.props.data.name}`)
        .then(data => {
            if (data.status !== 200 && data.status !== 401) toast.error(`Failed to delete file: (${data.status}) ${data.reason}`);
            else this.props.refresh();
        })

    move = () => {
        // Ensure path formatted correctly
        if (this.state.newPath[0] === "/") this.setState({ loading: true, newPath: "/" + this.state.newPath });
        else this.setState({ loading: true });

        Files.update(`${this.props.currentDirectory}/${this.props.data.name}`, "", this.state.newPath)
            .then(data => {
                if (data.status !== 200 && data.status !== 401) toast.error(`Failed to move file/directory: (${data.status}) ${data.reason}`);
                else this.props.refresh();
            })
            .finally(() => this.setState({ loading: false, renameModalOpen: false }));
    }
    rename = () => {
        this.setState({ loading: true });
        Files.update(`${this.props.currentDirectory}/${this.props.data.name}`, this.state.newName, "")
            .then(data => {
                if (data.status !== 200 && data.status !== 401) toast.error(`Failed to rename file/directory: (${data.status}) ${data.reason}`);
                else this.props.refresh()
            })
            .finally(() => this.setState({ loading: false, moveModalOpen: false }));
    }
    share = () => {
        this.setState({ loading: true });
        Shares.create(`${this.props.currentDirectory}/${this.props.data.name}`, this.state.shareTo)
            .then(data => {
                if (data.status !== 200 && data.status !== 401) toast.error(`Failed to share file: (${data.status}) ${data.reason}`);
            })
            .finally(() => this.setState({ loading: false, shareModalOpen: false }));
    }

    render() {
        let { data, onClick, currentDirectory } = this.props;
        return (
            <>
                <div className="card" style={{ marginTop: "0.25rem", marginBottom: "0.25rem"}}>
                    <div className="card-body" style={{ padding: "0.75rem"}}>
                        <div className="container">
                            <div className="row">
                                <div className="col-sm" style={{ paddingTop: "0.125rem" }} onClick={ data.directory ? onClick : () => {} }>
                                    <FontAwesomeIcon style={{ fontSize: "0.75rem" }} icon={ data.directory ? faFolder : faFileAlt }/> &nbsp;{data.name}
                                </div>
                                <div className="col-sm text-right">
                                    { !data.directory && (
                                        <>
                                            <button type="button" className="btn btn-outline-success btn-sm" style={{ fontSize: "0.75rem", marginRight: "0.25rem" }}
                                                    title="Download" onClick={this.download}><FontAwesomeIcon icon={faFileDownload}/></button>
                                            <button type="button" className="btn btn-outline-primary btn-sm" style={{ fontSize: "0.75rem", marginRight: "0.25rem" }}
                                                    title="Share" onClick={this.toggleShareModal.bind(this)}><FontAwesomeIcon icon={faShareAlt}/></button>
                                        </>
                                    )}
                                    <div className="btn-group">
                                        <button type="button" id="dropdownToggle" className="btn btn-outline-dark btn-sm dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false" style={{ fontSize: "0.75rem" }}
                                                title="Delete"><FontAwesomeIcon icon={faEllipsisV}/></button>
                                        <div className="dropdown-menu" aria-labelledby="dropdownToggle">
                                            <button type="button" className="dropdown-item" onClick={this.toggleRenameModal.bind(this)}>Rename</button>
                                            <button type="button" className="dropdown-item" onClick={this.toggleMoveModal.bind(this)}>Move</button>
                                            <button type="button" className="dropdown-item text-danger" onClick={this.remove}>Delete</button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <ReactModal isOpen={this.state.renameModalOpen} contentLabel="Rename Modal">
                    <div className="container">
                        <div className="row">
                            <div className="col-sm text-center">
                                <h5>Rename <code>{data.name}</code></h5>
                            </div>
                        </div>
                        <hr/>
                        <br/>
                        <div className="row">
                            <p>Current Name: <code>{data.name}</code></p>
                            <div className="input-group">
                                <div className="input-group-prepend">
                                    <span className="input-group-text">New Filename</span>
                                </div>
                                <input type="text" aria-label="New Filename" className="form-control" value={this.state.newName} onChange={this.onNameChange.bind(this)}/>
                            </div>
                        </div>
                        <br/>
                        <hr/>
                        <div className="row">
                            <div className="col-sm text-right">
                                { this.state.loading && (
                                    <div className="spinner-border" role="status">
                                        <span className="sr-only">Loading...</span>
                                    </div>
                                )}
                                { !this.state.loading && (
                                    <>
                                        <button type="button" className="btn btn-outline-danger" style={{ marginRight: "0.5rem"}} onClick={this.toggleRenameModal.bind(this)}>Cancel</button>
                                        <button type="button" className="btn btn-primary" onClick={this.rename.bind(this)}>Rename</button>
                                    </>
                                )}
                            </div>
                        </div>
                    </div>
                </ReactModal>
                <ReactModal isOpen={this.state.moveModalOpen} contentLabel="Move Modal">
                    <div className="container">
                        <div className="row">
                            <div className="col-sm text-center">
                                <h5>Move <code>{data.name}</code></h5>
                            </div>
                        </div>
                        <hr/>
                        <br/>
                        <div className="row">
                            <p>Current Path: <code>{currentDirectory}/{data.name}</code></p>
                            <div className="input-group">
                                <div className="input-group-prepend">
                                    <span className="input-group-text">New Path</span>
                                    <span className="input-group-text">/</span>
                                </div>
                                <input type="text" aria-label="New Filename" className="form-control" value={this.state.newPath} onChange={this.onPathChange.bind(this)}/>
                            </div>
                        </div>
                        <br/>
                        <hr/>
                        <div className="row">
                            <div className="col-sm text-right">
                                { this.state.loading && (
                                    <div className="spinner-border" role="status">
                                        <span className="sr-only">Loading...</span>
                                    </div>
                                )}
                                { !this.state.loading && (
                                    <>
                                        <button type="button" className="btn btn-outline-danger" style={{ marginRight: "0.5rem"}} onClick={this.toggleMoveModal.bind(this)}>Cancel</button>
                                        <button type="button" className="btn btn-primary" onClick={this.move.bind(this)}>Move</button>
                                    </>
                                )}
                            </div>
                        </div>
                    </div>
                </ReactModal>
                <ReactModal isOpen={this.state.shareModalOpen} contentLabel="Share Modal">
                    <div className="container">
                        <div className="row">
                            <div className="col-sm text-center">
                                <h5>Share <code>{data.name}</code></h5>
                            </div>
                        </div>
                        <hr/>
                        <br/>
                        <div className="row">
                            <div className="input-group">
                                <div className="input-group-prepend">
                                    <span className="input-group-text">Share to</span>
                                </div>
                                <input type="text" aria-label="Share to" className="form-control" value={this.state.shareTo} onChange={this.onUserChange.bind(this)}/>
                            </div>
                        </div>
                        <br/>
                        <hr/>
                        <div className="row">
                            <div className="col-sm text-right">
                                { this.state.loading && (
                                    <div className="spinner-border" role="status">
                                        <span className="sr-only">Loading...</span>
                                    </div>
                                )}
                                { !this.state.loading && (
                                    <>
                                        <button type="button" className="btn btn-outline-danger" style={{ marginRight: "0.5rem"}} onClick={this.toggleShareModal.bind(this)}>Cancel</button>
                                        <button type="button" className="btn btn-primary" onClick={this.share.bind(this)}>Share</button>
                                    </>
                                )}
                            </div>
                        </div>
                    </div>
                </ReactModal>
            </>
        )
    }
}

Entry.propTypes = {
    data: PropTypes.shape({
        directory: PropTypes.bool,
        last_modified: PropTypes.number,
        name: PropTypes.string,
        size: PropTypes.number,
    }).isRequired,
    currentDirectory: PropTypes.string.isRequired,
    onClick: PropTypes.func.isRequired,
    refresh: PropTypes.func.isRequired
}
