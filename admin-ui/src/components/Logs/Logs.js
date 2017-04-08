import React, { Component } from 'react';
import { Item } from 'semantic-ui-react';
import apiClient from '../../helpers/apiClient';
import { parseLogs } from '../../helpers/logParser';

class Logs extends Component {
    state = {
        logs: []
    };

    componentDidMount() {
        this.loadLogs().then(logs => {
            this.setState({ logs });
        });
    }    

    loadLogs = () => 
        apiClient.get('/rest/log', { responseType: 'text' })
        .then(response => parseLogs(response.data));

    render() {
        return (
            <Item.Group>
            { 
                this.state.logs.map((log, index) => (
                    log.method ? (
                        <Item key={index}>
                            <Item.Content>
                            <Item.Header>{log.method} {log.path}</Item.Header>
                            <Item.Meta>{log.time}</Item.Meta>
                            <Item.Description>
                                { Object.keys(log).map(k => k + '=' + log[k]).join(' ') }
                            </Item.Description>
                            </Item.Content>
                        </Item>
                    ) : (
                        <Item key={index}>
                            <Item.Content>
                            <Item.Header>INFO</Item.Header>
                            <Item.Description>
                                { log.info }
                            </Item.Description>
                            </Item.Content>
                        </Item>
                    )
                ))
            }
            </Item.Group>
        );
    }
}

export default Logs;
