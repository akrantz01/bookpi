import React, { Component } from 'react';
import { Container, Row, Col } from 'react-bootstrap';

class Footer extends Component {
    render() {
        return (
            <footer style={{ position: "absolute", bottom: 0, width: "100%", lineHeight: "45px", backgroundColor: "#f5f5f5"}}>
                <Container>
                    <Row className="justify-content-between">
                        <Col md>
                            <p className="text-muted" style={{ marginBottom: 0 }}>BookPi &copy; <a style={{ color: "darkgray" }} href="https://krantz.dev" target="_blank" rel="noopener noreferrer">Alexander Krantz</a> 2019</p>
                        </Col>
                    </Row>
                </Container>
            </footer>
        )
    }
}

export default Footer;
