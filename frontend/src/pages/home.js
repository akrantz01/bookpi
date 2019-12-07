import React, { Component } from 'react';
import { Container, Jumbotron, CardDeck, Card } from 'react-bootstrap';

class Home extends Component {
    render() {
        return (
            <Container>
                <Jumbotron>
                    <h1 className="display-4">BookPi</h1>
                    <hr className="my-4"/>
                    <p className="lead">Offline messaging, file storage, and file sharing</p>
                </Jumbotron>

                <CardDeck>
                    <Card>
                        <Card.Body>
                            <Card.Title><h4>Messaging</h4></Card.Title>
                            <Card.Subtitle><h6 className="mb-2 text-muted">Send messages to other users</h6></Card.Subtitle>
                            <Card.Text>Communicate with other people in realtime without being connected to the Internet. All communications are encrypted when they are stored in a database so they are not able to be read if the server is compromised. You can send normal text messages, images, and share files; group chats are also supported.</Card.Text>
                        </Card.Body>
                    </Card>

                    <Card>
                        <Card.Body>
                            <Card.Title><h4>Storage</h4></Card.Title>
                            <Card.Subtitle><h6 className="mb-2 text-muted">Securely store your files</h6></Card.Subtitle>
                            <Card.Text>Offline and secure storage for your files. All files are stored without their original names and encrypted to ensure that they are secure</Card.Text>
                        </Card.Body>
                    </Card>

                    <Card>
                        <Card.Body>
                            <Card.Title><h4>File Sharing</h4></Card.Title>
                            <Card.Subtitle><h6 className="mb-2 text-muted">Allow others to access your files</h6></Card.Subtitle>
                            <Card.Text>You can share your files with any other registered user. You also have the option to generate a link that anyone can access your files from. Not only can you share single files, but also folders.</Card.Text>
                        </Card.Body>
                    </Card>
                </CardDeck>
            </Container>
        )
    }
}

export default Home;
