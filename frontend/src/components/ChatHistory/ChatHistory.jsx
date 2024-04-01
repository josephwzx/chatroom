import React, { Component } from "react";
import "./ChatHistory.css";
import Message from "../Message";

class ChatHistory extends Component {
    render() {
        const messages = this.props.chatHistory.map((msg, index) => (
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
