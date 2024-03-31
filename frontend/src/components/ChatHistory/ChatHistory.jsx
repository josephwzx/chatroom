import React, { Component } from "react";
import "./ChatHistory.scss";
import Message from "../Message";

class ChatHistory extends Component {
    render() {
        const messages = this.props.chatHistory.map((msg, index) => (
            // Use the Message component for each message
            <Message key={index} message={msg} />
        ));

        return (
            <div className="ChatHistory">
                <h2>Chat History</h2>
                {messages}
            </div>
        );
    }
}

export default ChatHistory;
