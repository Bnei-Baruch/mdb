import React, { Component } from 'react';
import filesize from 'filesize';
import { Table } from 'semantic-ui-react';


const ObjectTable = ({ object }) => {
    if (!object) {
        return null;
    }

    return (
        <Table celled striped>
            <Table.Body>
                {
                    Object.keys(object).map(key => 
                        <Table.Row key={key}>
                            <Table.Cell collapsing>
                                <div>{ key }</div>
                            </Table.Cell>
                            <Table.Cell>
                                <div>{ object[key] }</div>
                            </Table.Cell>
                        </Table.Row>
                    )
                }
            </Table.Body>
        </Table>
    );
};

const transformValues = (key, value) => {
    switch (key) {
        case 'size':
            return filesize(value);
        default:
            return value;
    }
}

class File extends Component {

    state = {
        file: null
    };

    componentDidMount() {
        fetch(`http://rt-dev.kbb1.com:8080/admin/rest/files/${this.props.match.params.id}`)
            .then(response => {
                if (!response.ok) {
                    throw Error('Error loading files, response not ok.');
                }

                return response;
            })
            .then(response => response.json())
            .then(response => {
                this.setState({
                    file: response.file
                });
            });
    }

    render() {
        const { file } = this.state;
        if (!file) {
            return null;
        }

        return (
            <Table celled striped>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell colSpan='2'>File info</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    {
                        Object.keys(file).map(key => 
                            <Table.Row key={key}>
                                <Table.Cell collapsing>
                                    <div>{ key }</div>
                                </Table.Cell>
                                <Table.Cell>
                                    {
                                        key === 'properties' 
                                            ? <ObjectTable object={file[key]} />
                                            : <div>{ transformValues(key, file[key]) }</div>
                                    }
                                </Table.Cell>
                            </Table.Row>
                        )
                    }
                </Table.Body> 
            </Table>
        );
    }

}

export default File;
