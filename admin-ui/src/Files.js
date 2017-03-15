import React, { Component } from 'react';

class Files extends Component {
    constructor(props) {
        super(props);

        this.state = {
            files: [],
            showRemoveIcon: false,
            searchText: '',
        }

        this.getFiles('');
    }

    handleSearchChange = (e) => {
        const searchText = e.target.value;

        this.setState({
            searchText,
        });

        if (searchText === '') {
            this.setState({
                files: [],
                showRemoveIcon: false,
            });
        } else {
            this.setState({
                showRemoveIcon: true,
            });

            this.getFiles(searchText);
        }
    }

    getFiles = (searchText) => {
        return fetch('http://rt-dev.kbb1.com:8080/admin/rest/files')
        .then((response) => {
            return response.json().then(json => {
                if (json.status && json.status === 'ok') {
                    this.setState({files: json.files});
                } else {
                    console.error('Error fetching files');
                }
            });
        })
    }

    handleSearchCancel = () => {
        this.setState({
            files: [],
            showRemoveIcon: false,
            searchText: '',
        });
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
