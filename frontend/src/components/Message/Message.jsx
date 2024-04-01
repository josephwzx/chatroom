import React, { Component } from "react";
import "./Message.css";

class Message extends Component {
    render() {
        const { content, sender, created_at } = this.props.message;
        return (
            <div className="Message">
                <div className="Message-header">
                    {sender && <span className="sender"> {sender}</span>}
                    <span className="timestamp">{new Date(created_at).toLocaleString()}</span>
                </div>
                <p className="content">{content}</p>
            </div>
        );
    }
}

export default Message;
