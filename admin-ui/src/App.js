import React, { Component } from 'react';
import './App.css';
import Logs from './Logs.js';
import Files from './Files.js';

import { BrowserRouter as Router, Route } from 'react-router-dom'

class App extends Component {
  render() {
    return (
        <Router>
          <div>
            <Route exact path="/logs" component={Logs}/>
            <Route exact path="/files" component={Files}/>
          </div>
        </Router>
    );
  }
}

export default App;
