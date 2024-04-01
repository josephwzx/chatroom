import React, { Component } from "react";
import "./Message.css";

class Message extends Component {

    constructor(props) {
        super(props);
        this.state = {
            upvotecount: props.message.upvotecount,
            downvotecount: props.message.downvotecount,
        };
    }

    componentDidUpdate(prevProps) {
        if (this.props.message.upvotecount !== prevProps.message.upvotecount || this.props.message.downvotecount !== prevProps.message.downvotecount) {
            this.setState({
                upvotecount: this.props.message.upvotecount,
                downvotecount: this.props.message.downvotecount,
            });
        }
    }

    render() {
        const { id, content, sender, created_at, upvotecount, downvotecount} = this.props.message;
        
        return (
            <div className="Message">
                <div className="Message-header">
                    {sender && <span className="sender">{sender}</span>}
                    <span className="timestamp">{new Date(created_at).toLocaleString()}</span>
                </div>
                <p className="content">{content}</p>
                <div className="vote-buttons">
                    <button onClick={() => this.handleVote('up', id)} className="vote-button upvote">Upvote {this.state.upvotecount === null ? upvotecount : this.state.upvotecount}</button>
                    <button onClick={() => this.handleVote('down', id)} className="vote-button downvote">Downvote {this.state.downvotecount === null ? downvotecount : this.state.downvotecount}</button>
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
            this.setState({
                upvotecount: data.upvote_count,
                downvotecount: data.downvote_count,
            });
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
    }
}

export default Message;
