import React, { Component } from 'react';

class Files extends Component {
    constructor(props) {
        super(props);

        this.state = {
            files: [],
            showRemoveIcon: false,
            searchText: '',
        }

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
        Promise.resolve([{name: '1'}, {name: '2'}]).then((files) => {
            this.setState({
                files,
            });
        });
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
                <td>{'desc'}</td>
                <td className='right aligned'>{'some'}</td>
                <td className='right aligned'>{'info'}</td>
                <td className='right aligned'>{'fu'}</td>
                <td className='right aligned'>{'c'}</td>
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
                            <th className='eight wide'>Description</th>
                            <th>Kcal</th>
                            <th>Protein (g)</th>
                            <th>Fat (g)</th>
                            <th>Carbs (g)</th>
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
