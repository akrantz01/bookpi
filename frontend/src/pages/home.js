import React, { Component } from 'react'

class Home extends Component {
  render () {
    return (
      <div className="container">
        <div className="jumbotron">
          <h1 className="display-4">BookPi</h1>
          <hr className="my-4"/>
          <p className="lead">Offline messaging, file storage, and file sharing</p>
        </div>

        <div className="card-deck">
          <div className="card">
            <div className="card-body">
              <h4 className="card-title">Messaging</h4>
              <h6 className="card-subtitle mb-2 text-muted">Send messages to other users</h6>
              <p className="card-text">Communicate with other people in realtime without being connected
                                to the Internet. All communications are encrypted when they are stored in a database so
                                they are not able to be read if the server is compromised. You can send normal text
                                messages, images, and share files; group chats are also supported.</p>
            </div>
          </div>

          <div className="card">
            <div className="card-body">
              <h5 className="card-title">Storage</h5>
              <h6 className="card-subtitle mb-2 text-muted">Securely store your files</h6>
              <p className="card-text">Offline and secure storage for your files. All files are stored
                                without their original names and encrypted to ensure that they are secure.</p>
            </div>
          </div>

          <div className="card">
            <div className="card-body">
              <h5 className="card-title">File Sharing</h5>
              <h6 className="card-subtitle mb-2 text-muted">Allow others access to your files</h6>
              <p className="card-text">You can share your files with any other registered user. You also
                                have the option to generate a link that anyone can access your files from. Not only can
                                you share single files, but also folders.</p>
            </div>
          </div>
        </div>
      </div>
    )
  }
}

export default Home
