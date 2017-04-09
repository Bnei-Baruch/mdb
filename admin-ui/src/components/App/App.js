import React, { Component } from 'react';
import './App.css';
import Logs from '../Logs/Logs.js';
import Files from '../Files/Files.js';
import File from '../File/File.js';
import Welcome from '../Welcome/Welcome.js';
import { Button, Icon, Menu } from 'semantic-ui-react'
import { Router, NavLink, Route } from 'react-router-dom'

import createBrowserHistory from 'history/createBrowserHistory'

const history = createBrowserHistory()

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            activeItemsVisible: false,
            activeItems: [],
        };
        history.listen(this.historyChanged)
    }

    componentDidMount() {
        this.historyChanged(history.location);
    }

    historyChanged = (location) => {
        if (!!location.pathname.match(/files\/\d+/) &&
            !this.state.activeItems.includes(location.pathname)) {
            this.setState({
                activeItems: [
                    ...this.state.activeItems,
                    location.pathname
                ],
                activeItemsVisible: true,
            });
        }
    };

    activeItemText = (item) => {
        return item.match(/files\/(\d+)/)[1];
    };

    removeActiveItem = (item) => {
        if (this.state.activeItems.includes(item)) {
            const newActiveItems = this.state.activeItems.slice();
            newActiveItems.splice(newActiveItems.indexOf(item), 1);
            this.setState({
                activeItems: newActiveItems,
                activeItemsVisible: !!newActiveItems.length,
            });
        }
    }

    toggleActiveItems = () => this.setState({ activeItemsVisible: !this.state.activeItemsVisible });

    render() {
        const { activeItemsVisible } = this.state;
        return (
            <Router history={history}>
                <div style={{display: 'flex', flexDirection: 'column', height: '100vh'}}>
                    <Menu pointing>
                      <Menu.Item as={NavLink} to="/" exact>Welcome</Menu.Item>
                      <Menu.Item as={NavLink} to="/logs">Logs</Menu.Item>
                      <Menu.Item as={NavLink} to="/files">Files</Menu.Item>
                      <Menu.Menu position='right'>
                          <Button icon size='mini'
                                  style={{margin: 5}}
                                  onClick={this.toggleActiveItems}>
                            <Icon name='history' />
                          </Button>
                      </Menu.Menu>
                    </Menu>
                    <div style={{display: 'flex', flexDirection: 'row', flex: '1 0 auto'}}>
                        <div style={{display: 'flex', flexDirection: 'column', flex: '1 0 auto'}}>
                            <Route exact path="/" component={Welcome}/>
                            <Route exact path="/logs" component={Logs}/>
                            <Route exact path="/files" component={Files}/>
                            <Route exact path="/files/:id" component={File}/>
                        </div>
                        <div style={{
                            display: activeItemsVisible ? 'block' : 'none',
                            width: '150px'
                        }}>
                             <Menu fluid vertical tabular='right'>
                                 {
                                     this.state.activeItems.map(i => 
                                        <Menu.Item as={NavLink} key={i} to={i}>
                                            File #{this.activeItemText(i)}
                                            <i 
                                                className='remove icon'
                                                onClick={() => this.removeActiveItem(i)}
                                            />
                                        </Menu.Item>
                                 )}
                             </Menu>
                        </div>
                    </div>
                </div>
            </Router>
        );
    }
}

export default App;
