import React, { Component } from 'react';
import './App.css';
import Logs from './Logs.js';
import Files from './Files.js';
import Welcome from './Welcome.js';
import { Menu } from 'semantic-ui-react'
import { BrowserRouter, NavLink, Route } from 'react-router-dom'

import createBrowserHistory from 'history/createBrowserHistory'

const history = createBrowserHistory()

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            activeItem: '/',
        };
    }

    pathname = () => {};

    render() {
        return (
            <div>
                <BrowserRouter history={history}>
                    <div>
                        <Menu pointing>
                          <Menu.Item as={NavLink} to="/" exact>Welcome</Menu.Item>
                          <Menu.Item as={NavLink} to="/logs">Logs</Menu.Item>
                          <Menu.Item as={NavLink} to="/files">Files</Menu.Item>
                        </Menu>
                        <Route exact path="/" component={Welcome}/>
                        <Route exact path="/logs" component={Logs}/>
                        <Route exact path="/files" component={Files}/>
                    </div>
                </BrowserRouter>
            </div>
        );
    }
}

export default App;
