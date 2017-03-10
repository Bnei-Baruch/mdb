import React, { Component } from 'react';
import './App.css';
import Logs from './Logs.js';

import { BrowserRouter as Router, Route } from 'react-router-dom'

class App extends Component {
  render() {
    return (
        <Router>
          <div>
            <Route exact path="/admin/" component={Logs}/>
            <Route path="/admin/some" component={Logs}/>
          </div>
        </Router>
    );
  }
}

export default App;
