import React, { Component } from 'react';
import Spinner from './Spinner';

class Files extends Component {
    constructor(props) {
        super(props);

        this.state = {
            files: [],
            showRemoveIcon: false,
            searchText: '',
            loadingFiles: false,
            error: '',
        }
    }

    componentDidMount = () => {
        this.getFiles('');
    }

    handleSearchChange = (e) => {
        this.getFiles(e.target.value);
    }

    handleSearchCancel = () => {
        this.getFiles('');
    }

    getFiles = (searchText) => {
        this.setState({
            loadingFiles: true,
            searchText,
            showRemoveIcon: searchText !== '',
        });
        return fetch('http://rt-dev.kbb1.com:8080/admin/rest/files?query=' + searchText)
        .then((response) => {
            if (!response.ok) {
                throw Error('Error loading files, response not ok.');
            }
            this.setState({loadingFiles: false});
            return response.json().then(json => {
                if (json.status && json.status === 'ok') {
                    this.setState({files: json.files});
                } else {
                    throw Error('Error loading files, got bad status.');
                }
            });
        }).catch((e) => {
            this.setState({
                loadingFiles: false,
                error: 'Error loading files: ' + e
            });
        })
    }

    render() {
        const { showRemoveIcon, files } = this.state;
        const removeIconStyle = showRemoveIcon ? {} : { visibility: 'hidden' };

        const fileRows = files.map((file, idx) => (
            <tr key={idx}>
                <td>{file.uid}</td>
                <td>{file.name}</td>
                <td className='right aligned'>{file.file_created_at}</td>
            </tr>
        ));

        return (
            <div id='file-search'>
                <table className='ui selectable structured large table'>
                    <thead>
                        <tr>
                            <th colSpan='5'>
                                <div className='ui fluid search'>
                                <div className='ui icon input'>
                                    <input
                                        className='prompt'
                                        type='text'
                                        placeholder='Search files...'
                                        value={this.state.searchText}
                                        onChange={this.handleSearchChange}
                                    />
                                    <i className='search icon' />
                                </div>
                                <i
                                    className='remove icon'
                                    onClick={this.handleSearchCancel}
                                    style={removeIconStyle}
                                />
                                {this.state.loadingFiles ? <span><Spinner /> Searching...</span> : null}
                                {!!this.state.error ? <span style={{color: 'red'}}>{this.state.error}</span> : null}
                                </div>
                            </th>
                        </tr>
                        <tr>
                            <th>UID</th>
                            <th className='eight wide'>Name</th>
                            <th>created at</th>
                        </tr>
                    </thead>
                <tbody>
                {fileRows}
                </tbody>
                </table>
            </div>
        );
    }

}

export default Files;
