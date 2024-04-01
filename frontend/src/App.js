import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg } from "./api";
import Header from "./components/Header";
import ChatHistory from './components/ChatHistory/ChatHistory';
import Login from './components/Login/Login';
import Register from './components/Register/Register';
import ChatInput from './components/Input/Input';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            chatHistory: [],
            isAuthenticated: false,
            isRegistering: false,
            token: '',
        };
    }
    
    componentDidMount() {
        connect((msg) => {
            const message = JSON.parse(msg.data);
            const formedMessage = {
                content: message.content,
                sender: message.sender,
                created_at: message.created_at,
            }
            console.log("formed", formedMessage);

            this.setState(prevState => ({
                chatHistory: [...prevState.chatHistory, formedMessage]
            }), this.scrollToBottom);
        });
    }

    send(event) {
        if(event.keyCode === 13) {
            const username = localStorage.getItem('username');
            sendMsg(event.target.value, username);
            event.target.value = "";
        }
    }

    toggleRegister = () => {
        this.setState(prevState => ({
            isRegistering: !prevState.isRegistering
        }));
    };   

    // Adjusted to fit the new Login component
    handleLoginSuccess = (token) => {
        localStorage.setItem('token', token); // Store the token
        this.setState({ isAuthenticated: true }, () => {
            this.fetchChatHistory(); // Fetch chat history after successful login
        });
    };

    fetchChatHistory = () => {
        // Assuming the fetch URL and method are correct
        fetch('http://localhost:8080/history')
        .then(response => response.json())
        .then(history => this.setState({ chatHistory: history }))
        .catch(error => console.error('Error fetching chat history:', error));
    };

    render() {
        if (!this.state.isAuthenticated) {
            if (this.state.isRegistering) {
                <Header />
                return <Register onRegister={this.toggleRegister} />;
            }
            return (
                <div>
                    <Header />
                    <Login onLoginSuccess={this.handleLoginSuccess} />
                    <button onClick={this.toggleRegister}>Register</button>
                </div>
            );
        }

        return (
            <div className="App">
                <Header />
                <ChatHistory 
                    chatHistory={this.state.chatHistory}
                    />
                <ChatInput send={this.send} />
            </div>
        );
    }
}

export default App;
