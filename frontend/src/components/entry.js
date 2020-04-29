import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFileDownload, faFolder, faFileAlt, faEllipsisV, faShareAlt } from "@fortawesome/free-solid-svg-icons";

export default class Entry extends Component {
    render() {
        let { data, onClick } = this.props;
        return (
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
                                                title="Download"><FontAwesomeIcon icon={faFileDownload}/></button>
                                        <button type="button" className="btn btn-outline-primary btn-sm" style={{ fontSize: "0.75rem", marginRight: "0.25rem" }}
                                                title="Download"><FontAwesomeIcon icon={faShareAlt}/></button>
                                    </>
                                )}
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

Entry.propTypes = {
    data: PropTypes.shape({
        directory: PropTypes.bool,
        last_modified: PropTypes.number,
        name: PropTypes.string,
        size: PropTypes.number,
    }).isRequired,
    onClick: PropTypes.func.isRequired
}
