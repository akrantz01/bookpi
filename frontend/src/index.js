import React from 'react'
import ReactDOM from 'react-dom'
import ReactModal from 'react-modal'
import * as serviceWorker from './serviceWorker'
import 'bootstrap/dist/js/bootstrap.min'
import 'bootstrap/dist/css/bootstrap.min.css'
import 'popper.js/dist/popper.min'
import 'jquery/dist/jquery.min'
import 'react-toastify/dist/ReactToastify.min.css'

import Router from './router'

ReactModal.setAppElement('#root')
ReactModal.defaultStyles = {
  overlay: {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(255, 255, 255, 0.75)'
  },
  content: {
    position: 'absolute',
    top: '25%',
    left: '25%',
    width: '50%',
    background: '#fff',
    border: '1px solid #ccc',
    overflow: 'auto',
    WebkitOverflowScrolling: 'touch',
    borderRadius: '8px',
    outline: 'none',
    padding: '20px'
  }
}

ReactDOM.render(<Router/>, document.getElementById('root'))

serviceWorker.unregister()
