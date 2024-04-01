import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg } from "./api";
import Header from "./components/Header";
import ChatHistory from './components/ChatHistory/ChatHistory';
import Login from './components/Login/Login'; // Ensure this is the new Login component
import Register from './components/Register/Register';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            chatHistory: [],
            isAuthenticated: false,
            isRegistering: false,
            token: '', // Store token here if needed
        };
    }

    componentDidMount() {
        connect((msg) => {
            console.log("New Message");
            this.setState(prevState => ({
                chatHistory: [...prevState.chatHistory, msg]
            }));
        });
    }

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
                return <Register onRegister={this.toggleRegister} />;
            }
            return (
                <div>
                    <Login onLoginSuccess={this.handleLoginSuccess} />
                    <button onClick={this.toggleRegister}>Register</button>
                </div>
            );
        }

        return (
            <div className="App">
                <Header />
                <ChatHistory chatHistory={this.state.chatHistory} />
                <button onClick={() => sendMsg("hello")}>Hit</button>
            </div>
        );
    }
}

export default App;
