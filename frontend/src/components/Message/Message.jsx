import React, { Component } from "react";
import "./Message.css";

class Message extends Component {
    render() {
        // Assuming the message prop is now an object instead of a JSON string
        const { content, sender, created_at } = this.props.message;
        return (
            <div className="Message">
                <p>{content}</p>
                {sender && <p className="sender">Sender: {sender}</p>}
                <p className="timestamp">Sent: {new Date(created_at).toLocaleString()}</p>
            </div>
        );
    }
}

export default Message;
