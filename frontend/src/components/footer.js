import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { withRouter } from 'react-router-dom'

class Footer extends Component {
  render () {
    if (this.props.location.pathname === '/sign-in' || this.props.location.pathname === '/sign-up') return (<></>)

    return (
      <footer style={{ position: 'fixed', bottom: 0, left: 0, width: '100%', lineHeight: '45px', backgroundColor: '#f5f5f5' }}>
        <div className="container">
          <div className="row justify-content-between">
            <div className="col-md">
              <p className="text-muted" style={{ marginBottom: 0 }}>BookPi &copy; <a style={{ color: 'darkgray' }} href="https://krantz.dev" target="_blank" rel="noopener noreferrer">Alexander Krantz</a> 2020</p>
            </div>
          </div>
        </div>
      </footer>
    )
  }
}

Footer.propTypes = {
  location: PropTypes.object
}

export default withRouter(Footer)
