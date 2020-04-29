import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFileDownload, faFolder, faFileAlt, faEllipsisV, faShareAlt } from "@fortawesome/free-solid-svg-icons";

export class FileEntry extends Component {
    render() {
        let { file } = this.props;
        return (
            <div className="card" style={{ marginTop: "0.25rem", marginBottom: "0.25rem"}}>
                <div className="card-body" style={{ padding: "0.75rem"}}>
                    <div className="container">
                        <div className="row">
                            <div className="col-sm" style={{ paddingTop: "0.125rem" }}>
                                <FontAwesomeIcon style={{ fontSize: "0.75rem" }} icon={faFileAlt}/> &nbsp;{file.name}
                            </div>
                            <div className="col-sm text-right">
                                <button type="button" className="btn btn-outline-success btn-sm" style={{ fontSize: "0.75rem", marginRight: "0.25rem" }}
                                        title="Download"><FontAwesomeIcon icon={faFileDownload}/></button>
                                <button type="button" className="btn btn-outline-primary btn-sm" style={{ fontSize: "0.75rem", marginRight: "0.25rem" }}
                                        title="Download"><FontAwesomeIcon icon={faShareAlt}/></button>
                                <div className="btn-group">
                                    <button type="button" id="dropdownToggle" className="btn btn-outline-dark btn-sm dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false" style={{ fontSize: "0.75rem" }}
                                            title="Delete"><FontAwesomeIcon icon={faEllipsisV}/></button>
                                    <div className="dropdown-menu" aria-labelledby="dropdownToggle">
                                        <button type="button" className="dropdown-item">Rename</button>
                                        <button type="button" className="dropdown-item">Delete</button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

FileEntry.propTypes = {
    file: PropTypes.shape({
        directory: PropTypes.bool,
        last_modified: PropTypes.number,
        name: PropTypes.string,
        size: PropTypes.number
    }).isRequired
};


export class DirectoryEntry extends Component {
    render() {
        let { directory, onClick } = this.props;
        return (
            <div className="card" style={{ marginTop: "0.25rem", marginBottom: "0.25rem" }}>
                <div className="card-body" style={{ padding: "0.75rem" }}>
                    <div className="container">
                        <div className="row">
                            <div className="col-sm" style={{ paddingTop: "0.125rem" }} onClick={onClick}>
                                <FontAwesomeIcon style={{ fontSize: "0.75rem" }} icon={faFolder}/> &nbsp;{directory.name}
                            </div>
                            <div className="col-sm text-right">
                                <div className="btn-group">
                                    <button type="button" id="dropdownToggle" className="btn btn-outline-dark btn-sm dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false" style={{ fontSize: "0.75rem" }}
                                            title="Delete"><FontAwesomeIcon icon={faEllipsisV}/></button>
                                    <div className="dropdown-menu" aria-labelledby="dropdownToggle">
                                        <button type="button" className="dropdown-item">Rename</button>
                                        <button type="button" className="dropdown-item">Delete</button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

DirectoryEntry.propTypes = {
    directory: PropTypes.shape({
        directory: PropTypes.bool,
        last_modified: PropTypes.number,
        name: PropTypes.string,
        size: PropTypes.number,
    }),
    onClick: PropTypes.func.isRequired
}
