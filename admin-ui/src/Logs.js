import React, { Component } from 'react';
import { Item } from 'semantic-ui-react';

class Logs extends Component {
    constructor() {
        super();

        this.state = { logs: [] };
    }

    componentDidMount() {
        this.loadLogs().then(logs => {
            this.setState({ logs });
        });
    }

    parseLogs(text) {
        console.log(text);
        const rawLogs = text.split('\n').reduce((a, line) => {
            if (a.length === 0) {
                a.push([line]);
            } else {
                const lastLog = a[a.length - 1];
                const lastLogLine = lastLog[lastLog.length - 1];
                if (this.isEndLogLine(lastLogLine)) {
                    a.push([line]);
                } else {
                    lastLog.push(line);
                }
            }
            return a;
        }, []);

        const logs = rawLogs.map(lines => {
            if (this.isEndLogLine(lines[lines.length - 1])) {
                const info = lines.slice(0, lines.length - 2).join('\n');
                const httpObj = this.parseLine(lines[lines.length - 1]);
                httpObj.info = [info, httpObj.info].filter(i => i).join('\n');
                return httpObj;
            } else {
                return { info: lines.join('\n') };
            }
        });

        console.log(logs);
        return logs;
    }

    parseLine(line) {
        if (line.indexOf('time=') !== -1) {
            return this.parseHTTP(line);
        } else {
            return { info: line };
        }
    }

    parseHTTP(line) {
        const idx = line.indexOf('=');
        if (idx === -1) {
            return {};
        } else {
            const name = line.substring(0, idx);
            let value = '';
            let startOfValue = idx + 1;
            let endOfValue = idx + 1;
            if (startOfValue < line.length - 1) {
                if (line[startOfValue] === '"') {
                    startOfValue++;
                    endOfValue = line.indexOf('"', idx + 2);
                    value = line.substring(startOfValue, endOfValue);
                    endOfValue++;
                } else {
                    endOfValue = line.indexOf(' ', idx + 1);
                    value = line.substring(startOfValue, endOfValue);
                }
            }
            const ret = this.parseHTTP(line.substring(endOfValue + 1));
            ret[name] = value;
            return ret;
        }
    }

    isEndLogLine(line) {
        return line.indexOf('time=') !== -1 ||
            line.indexOf('[GIN-debug]') !== -1;
    }

    loadLogs() {
        return fetch('http://rt-dev.kbb1.com:8080/admin/rest/log')
        .then((response) => {
            return response.text().then(text => {
                const logs = this.parseLogs(text);
                return logs;
            });
        })
    }

    render() {
        return (
            <Item.Group>
            { this.state.logs.map((log, index) => (
                    log.method ?
                    <Item key={index}>
                    <Item.Content>
                    <Item.Header>{log.method} {log.path}</Item.Header>
                    <Item.Meta>{log.time}</Item.Meta>
                    <Item.Description>
                    { Object.keys(log).map(k => k + '=' + log[k]).join(' ') }
                    </Item.Description>
                    </Item.Content>
                    </Item>
                    :
                    <Item key={index}>
                    <Item.Content>
                    <Item.Header>INFO</Item.Header>
                    <Item.Description>
                    { log.info }
                    </Item.Description>
                    </Item.Content>
                    </Item>
            )) }
            </Item.Group>
        );
    }
}

export default Logs
