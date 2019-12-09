import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Chats, Messages } from '../api';

async function getChats() {
    // Get list of user chats
    let chatList = await Chats.list();
    if (chatList.status !== 200) toast.error(`Failed loading chat list: (${chatList.status}) ${chatList.reason}`);

    let chatData = [];

    // Get specific data about each
    for (let id of chatList.data) {
        // Describe chat
        let chat = await Chats.read(id);
        if (chat.status !== 200) toast.error(`Failed to load chat data for ${id}: (${chat.status}) ${chat.reason}`);

        // Add id to chat
        chat.data.id = id;

        // Add to array
        chatData.push(chat.data);
    }

    return chatData;
}

class Chat extends Component {
    constructor(props) {
        super(props);

        this.state = {
            selected: {},
            chats: [],
            loadingList: true,
            loadingChat: false,
            message: "",
            initialMessage: "",
            to: ""
        }
    }

    componentDidMount() {
        getChats().then(chats => this.setState({ chats, loadingList: false }));
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (prevState.selected.id !== this.state.selected.id) {
            Messages.list(this.state.selected.id).then(data => {
                if (data.status !== 200) toast.error(`Failed to load messages for ${this.state.selected.id}: (${data.status}) ${data.reason}`);
                else this.setState({ selected: { ...this.state.selected, messages: data.data }, loadingChat: false });
            });
        }
    }

    refreshChats = () => {
        this.setState({ loadingList: true });
        getChats().then(chats => this.setState({ chats, loadingList: false }));
    };
    selectChat = event => (this.state.selected.id !== event.target.dataset.id && event.target.dataset.user1 && event.target.dataset.user2) ? this.setState( { selected: Object.assign({}, event.target.dataset), loadingChat: true }) : false;
    createChat = () => Chats.create(this.state.to, this.state.initialMessage).then(data => {
        if (data.status === 400 && data.reason === "Specified recipient does not exist") toast.error("Specified user does not exist");
        else if (data.status !== 200) toast.error(`Failed to create new chat: (${data.status}) ${data.reason}`);
        else {
            this.setState({ to: "", initialMessage: "" });
            this.refreshChats();
        }
    });
    deleteChat = event => Chats.delete(event.target.dataset.id).then(data => {
        if (data.status !== 200) toast.error(`Failed to delete chat ${event.target.dataset.id}: (${data.status}) ${data.reason}`);
        else this.refreshChats();
    });

    sendMessage = () => Messages.create(this.state.selected.id, this.state.message).then(data => {
        if (data.status !== 200) toast.error(`Failed to send message to ${this.state.selected.id}: (${data.status}) ${data.reason}`);
        else this.setState({ selected: { ...this.state.selected, messages: [...this.state.selected.messages, `${this.props.username}:${this.state.message}`] }, message: "" });
    });

    onMessageInput = event => this.setState({ message: event.target.value });
    onInitialMessageInput = event => this.setState({ initialMessage: event.target.value });
    onToInput = event => this.setState({ to: event.target.value });

    render() {
        if (!this.props.loggedIn) this.props.history.push("/sign-in");

        return (
            <div className="container">
                <div className="modal fade" id="createModal" tabIndex="-1" role="dialog" aria-labelledby="createModalLabel" aria-hidden="true">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">
                            <div className="modal-header">
                                <h5 className="modal-title" id="createModalLabel">Create a Chat</h5>
                                <button type="button" className="close" data-dismiss="modal" aria-label="Close">
                                    <span aria-hidden="true">&times;</span>
                                </button>
                            </div>
                            <div className="modal-body">
                                <form>
                                    <label htmlFor="to">To:</label>
                                    <input type="text" onChange={this.onToInput.bind(this)} value={this.state.to} className="form-control" id="to" autoFocus required/>

                                    <br/>

                                    <label htmlFor="initialMessage">Initial Message:</label>
                                    <input type="text" onChange={this.onInitialMessageInput.bind(this)} value={this.state.initialMessage} className="form-control" id="initialMessage" required/>

                                    <br/>

                                    <button style={{ marginRight: "10px" }} type="button" className="btn btn-secondary" data-dismiss="modal">Close</button>
                                    <button type="submit" className="btn btn-primary" data-dismiss="modal" onClick={this.createChat.bind(this)}>Create</button>
                                </form>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="row">
                    <div className="col-sm-4">
                        <div className="card" style={{ height: "85vh" }}>
                            <div className="card-body" style={{ overflow: "auto" }}>
                                <h3 className="card-title">Chats</h3>
                                <div className="btn-group" role="group" aria-label="Chat operations">
                                    <button type="button" className="btn btn-sm btn-outline-success" data-toggle="modal" data-target="#createModal">New</button>
                                    <button type="button" className="btn btn-sm btn-outline-primary" onClick={this.refreshChats.bind(this)}>Refresh</button>
                                </div>
                                <hr className="my-4"/>
                                <div className="list-group list-group-flush">
                                    { this.state.loadingList && (
                                        <div className="spinner-border" role="status">
                                            <span className="sr-only">Loading...</span>
                                        </div>
                                    )}
                                    { !this.state.loadingList && this.state.chats.map(chat => (
                                        <button type="button" className="list-group-item list-group-item-action" key={chat.id} data-id={chat.id} data-user1={chat.user1} data-user2={chat.user2} onClick={this.selectChat.bind(this)}>
                                            <span onClick={this.deleteChat.bind(this)} data-id={chat.id} aria-hidden="true">&times;</span>
                                            &nbsp;{ (this.props.username === chat.user1) ? chat.user2 : chat.user1 }
                                        </button>)) }
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="col-sm-8">
                        <div className="card" style={{ height: "85vh" }}>
                            <h5 className="card-header">
                                { !this.state.selected.id && "Select a chat to begin" }
                                { this.state.selected.id && `Chat with ${ (this.props.username === this.state.selected.user1) ? this.state.selected.user2 : this.state.selected.user1 }` }
                            </h5>
                            <div className="card-body" style={{ overflow: "auto" }}>
                                { this.state.loadingChat && (
                                    <div className="spinner-border" role="status">
                                        <span className="sr-only">Loading...</span>
                                    </div>
                                )}
                                { !this.state.loadingChat && this.state.selected.messages && this.state.selected.messages.map(message => {
                                    let user = message.split(":", 1)[0];
                                    let text = message.split(":", 2)[1];
                                    return <p key={message+Math.random().toString()}><span className={`badge badge-${(user === this.props.username) ? "primary" : "dark"}`}>{ (user === this.props.username) ? "You" : user }</span> { text }</p>
                                })}
                            </div>
                            <div className="card-footer">
                                <form className="form-inline">
                                    <label className="sr-only" htmlFor="message">Message</label>
                                    <input type="text" className="form-control mb-2 mr-sm-2 w-75" id="message" onInput={this.onMessageInput.bind(this)} value={this.state.message} disabled={this.state.loadingChat || !this.state.selected.messages }/>

                                    <button type="submit" className="btn btn-success mb-2" disabled={this.state.loadingChat || !this.state.selected.messages } onClick={this.sendMessage.bind(this)}>Send</button>
                                </form>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

Chat.propTypes = {
    loggedIn: PropTypes.bool,
    history: PropTypes.object,
    username: PropTypes.string,
};

export default withRouter(Chat);
