import React, { Component } from 'react';
import './App.css';
import Logs from './Logs.js';
import Files from './Files.js';
import File from './File.js';
import Welcome from './Welcome.js';
import { Button, Icon, Menu } from 'semantic-ui-react'
import { Router, NavLink, Route } from 'react-router-dom'

import createBrowserHistory from 'history/createBrowserHistory'

const history = createBrowserHistory()

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            historyVisible: false,
        };
    }

    toggleHistory = () => this.setState({ historyVisible: !this.state.historyVisible })

    render() {
        const { historyVisible } = this.state;
        return (
            <Router history={history}>
                <div>
                    <Menu pointing>
                      <Menu.Item as={NavLink} to="/" exact>Welcome</Menu.Item>
                      <Menu.Item as={NavLink} to="/logs">Logs</Menu.Item>
                      <Menu.Item as={NavLink} to="/files">Files</Menu.Item>
                      <Menu.Menu position='right'>
                          <Button icon size='mini'
                                  style={{margin: 5}}
                                  onClick={this.toggleHistory}>
                            <Icon name='history' />
                          </Button>
                      </Menu.Menu>
                    </Menu>
                    <div style={{display: 'flex', flexDirection: 'row'}}>
                        <div style={{display: historyVisible ? 'block' : 'none',
                                     height: '100vh',
                                     width: '200px',
                                     float: 'right'}}>
                            This is History!
                        </div>
                        <div>
                            <div style={{display: 'flex', flexDirection: 'column', height: '100vh'}}>
                                <Route exact path="/" component={Welcome}/>
                                <Route exact path="/logs" component={Logs}/>
                                <Route exact path="/files" component={Files}/>
                                <Route exact path="/files/:id" component={File}/>
                            </div>
                        </div>
                    </div>
                </div>
            </Router>
        );
    }
}

export default App;
