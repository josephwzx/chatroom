import React, { Component } from "react";
import "./Message.css";

class Message extends Component {
    render() {
        console.log(this.props.message)
        const { id, content, sender, created_at, upvote_count, downvote_count } = this.props.message;
        return (
            <div className="Message">
                <div className="Message-header">
                    {sender && <span className="sender">{sender}</span>}
                    <span className="timestamp">{new Date(created_at).toLocaleString()}</span>
                </div>
                <p className="content">{content}</p>
                <div className="vote-buttons">
                    <button onClick={() => this.handleVote('up', id)} className="vote-button upvote">Upvote {upvote_count}</button>
                    <button onClick={() => this.handleVote('down', id)} className="vote-button downvote">Downvote {downvote_count}</button>
                </div>
            </div>
        );
    }

    handleVote = (type, id) => {
        const vote = {
            message_id: id,
            vote_type: type === 'up' ? 'upvote' : 'downvote'
        };
        console.log(localStorage.getItem('token'))
        fetch('http://localhost:8080/vote', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify(vote),
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json(); // You might send back the updated vote count from the server
        })
        .then(data => {
            // Here you could update the state to reflect the new vote count if needed
            console.log('Vote successfully sent to server', data);
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
    }
}

export default Message;
