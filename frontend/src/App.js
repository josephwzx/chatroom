import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg } from "./api";
import Header from "./components/Header";
import ChatHistory from './components/ChatHistory/ChatHistory';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            chatHistory: []
        };
    }

    componentDidMount() {
        connect((msg) => {
            console.log("New Message");
            this.setState(prevState => ({
                chatHistory: [...prevState.chatHistory, msg]
            }));
            console.log(this.state);
        });
        fetch('http://localhost:8080/history')
        .then(response => {
          console.log(response); // Log the full response for debugging
          if (!response.ok) {
            throw new Error(`Network response was not ok, status: ${response.status}`);
          }
          const contentType = response.headers.get('Content-Type');
          if (!contentType || !contentType.includes('application/json')) {
            throw new Error(`Unexpected content type: ${contentType}`);
          }
          return response.json(); // Properly parse and return JSON data once
        })
            .then(history => {
                // Here we set the state with the history received
                this.setState({ chatHistory: history });
                console.log("History fetched:", history);
            })
            .catch(error => console.error('Error fetching chat history:', error));
    }

    send() {
        console.log("hello");
        sendMsg("hello");
    }

    render() {
        return (
            <div className="App">
                <Header />
                <ChatHistory chatHistory={this.state.chatHistory} /> {/* ChatHistory now receives the history as a prop */}
                <button onClick={this.send}>Hit</button>
            </div>
        );
    }
}

export default App;
