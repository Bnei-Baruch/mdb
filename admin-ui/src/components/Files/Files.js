import React, { Component } from 'react';
import Spinner from '../Spinner/Spinner';
import './Files.css';
import { AutoSizer, Column, InfiniteLoader, Table } from 'react-virtualized';
import 'react-virtualized/styles.css';
import { Link } from 'react-router-dom';
import apiClient from '../../helpers/apiClient';

const RowRenderer = ({ className, columns, key, style, index, rowData }) => {
    if(!rowData || !rowData.id) {
        return (
            <div 
                className={className}
                key={key}
                style={style}
            >
                 <div 
                    className="flex-center-center"
                    style={{height: '100%', width: '100%'}}
                >
                    <div className="placeholder" />
                 </div>
            </div>
        );
    }

    return (
        <div
            className={className}
            key={key}
            style={style}
        >
             {columns}
        </div>
    );
};

const LinkToFileCellRenderer = ({ cellData, dataKey }) =>
    <Link to={`/files/${cellData}`}>{cellData}</Link>;

class Files extends Component {
    constructor(props) {
        super(props);

        this.firstLimit = 100;
    }

    state = {
        // Should be eventually props.
        files: [],
        resetFiles: false,
        matching: 0,
        total: 0,

        loadingFiles: false,
        error: '',

        // Should be eventually state.
        showRemoveIcon: false,
        searchText: '',
    };

    componentDidMount = () => {
        this.searchFiles('', 0, this.firstLimit);
    }

    handleSearchChange = (e) => {
        this.searchFiles(e.target.value, 0, this.firstLimit);
    }

    handleSearchCancel = () => {
        this.searchFiles('', 0, this.firstLimit);
    }

    searchFiles = (searchText, startIndex, stopIndex) => {
        console.log('Search text:', searchText, 'Fetching start: ' + startIndex + ' stop: ' + stopIndex);
        const limit = stopIndex - startIndex + 1;
        this.setState((prevState) => {
            const newState = {
                loadingFiles: true,
                searchText,
                showRemoveIcon: searchText !== '',
                resetFiles: prevState.searchText !== searchText,
            };
            if (newState.resetFiles) {
                this.refs.inf.resetLoadMoreRowsCache();
            }
            return newState;
        }, () => {
            apiClient.get('/rest/files', {
                params: {
                    offset: startIndex,
                    limit,
                    query: searchText    
                }
            }).then(response => {
                const { files, matching, total } = response.data;
                this.setState((prevState) => {
                    const newFiles = prevState.resetFiles ? [] : prevState.files;
                    files.forEach((f, i) => {
                        newFiles[i + startIndex] = f;
                        f.index = i + startIndex;
                    });
                    return {
                        loadingFiles: false,
                        resetFiles: false,
                        files: newFiles,
                        matching,
                        total,
                        error: '',
                    };
                });
            }).catch((e) => {
                console.log(e);
                this.setState({
                    loadingFiles: false,
                    error: 'Error loading files: ' + e
                });
            })
        });
    }

    isRowLoaded = ({ index }) => {
        const item = this.state.files[index];
        return item && typeof item.id !== 'undefined';
    };

    rowGetter = ({ index }) => 
        this.isRowLoaded({ index }) ? this.state.files[index] : {};

    loadMoreRows = ({ startIndex, stopIndex }) => 
        this.searchFiles(this.state.searchText, startIndex, stopIndex);

    render() {
        const { 
            showRemoveIcon, 
            searchText, 
            loadingFiles, 
            error, 
            matching, 
            total 
        } = this.state;

        return (
            <div style={{ display: 'flex', flex: '1 1 auto', flexDirection: 'column'}}>
                <Header 
                    showRemoveIcon={showRemoveIcon}
                    searchText={searchText}
                    handleSearchChange={this.handleSearchChange}
                    handleSearchCancel={this.handleSearchCancel}
                    loadingFiles={loadingFiles}
                    error={error}
                    matching={matching}
                    total={total} 
                />
                <div style={{ display: 'flex', flex: '1 1 auto', flexDirection: 'column'}}>
                    <InfiniteLoader 
                        ref="inf"
                        isRowLoaded={this.isRowLoaded}
                        threshold={100}
                        loadMoreRows={this.loadMoreRows}
                        rowCount={matching}
                    >
                        {({ onRowsRendered, registerChild }) => (
                            <AutoSizer>
                                {({ width, height }) => (
                                    <Table 
                                        headerHeight={50}
                                        height={height}
                                        width={width}
                                        rowCount={matching}
                                        ref={registerChild}
                                        onRowsRendered={onRowsRendered}
                                        rowRenderer={RowRenderer}
                                        rowGetter={this.rowGetter}
                                        rowHeight={50}
                                    >
                                        <Column 
                                            label='Index'
                                            cellDataGetter={({ rowData }) => rowData.index}
                                            dataKey='index'
                                            width={60} 
                                        />
                                        <Column 
                                            label='ID'
                                            dataKey='id'
                                            cellRenderer={LinkToFileCellRenderer}
                                            width={80} 
                                        />
                                        <Column 
                                            label='UID'
                                            dataKey='uid'
                                            width={80} 
                                        />
                                        <Column 
                                            label='Name'
                                            dataKey='name'
                                            width={160} flexGrow={1}
                                        />
                                        <Column 
                                            label='Created at'
                                            dataKey='file_created_at'
                                            width={80}
                                            flexGrow={1} 
                                        />
                                    </Table>
                                )}
                            </AutoSizer>
                        )}
                    </InfiniteLoader>
                </div>
            </div>
        );
    }

}

const Header = (props) => {
    const removeIconStyle = props.showRemoveIcon ? {} : { visibility: 'hidden' };

    return (
        <div 
            className='ui fluid search'
            style={{ 
                display: 'flex',
                flexDirection: 'row',
                justifyContent: 'space-between',
                paddingLeft: 10,
                paddingRight: 10 
            }}
        >
            <div>
                <div className='ui icon input'>
                    <input 
                        className='prompt'
                        type='text'
                        placeholder='Search files...'
                        value={props.searchText}
                        onChange={props.handleSearchChange} 
                    />
                    <i className='search icon' />
                </div>
                <i 
                    className='remove icon'
                    onClick={props.handleSearchCancel}
                    style={removeIconStyle} 
                />
            </div>
            <div className='flex-space-between-center'>
                {
                    props.loadingFiles &&
                        <span className='flex-space-between-center'>
                            <Spinner/>
                            <span style={{marginLeft: '10px'}}>Searching...</span>
                        </span>
                }
                {!!props.error && <span style={{color: 'red', marginLeft: '10px'}}>{props.error}</span>}
            </div>
            <div className='flex-space-between-center'>
                {
                    props.matching >= 0 && props.total >= 0 && 
                        <span>Matched {props.matching} of {props.total}</span>
                    
                }
            </div>
        </div>
    );
}

export default Files;
