import React, { Component } from 'react';
import Spinner from './Spinner';
import './Files.css';
import { AutoSizer, Column, InfiniteLoader, Table } from 'react-virtualized'
import 'react-virtualized/styles.css'
import { Link } from 'react-router-dom'

const RowRenderer = ({ className, columns, key, style, index, rowData }) => {
    if(!rowData || !rowData.id) {
        return (
            <div className={className}
                 key={key}
                 style={style}>
                 <div className='flex-center-center'
                      style={{height: '100%', width: '100%'}}>
                     <div className='placeholder'>
                     </div>
                 </div>
            </div>
        );
    }

    return (
        <div className={className}
             key={key}
             style={style} >
             {columns}
        </div>
    );
};

const LinkToFileCellRenderer = ({ cellData, dataKey }) =>
    <Link to={["/files", cellData].join('/')}>{cellData}</Link>;

class Files extends Component {
    constructor(props) {
        super(props);

        this.state = {
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
        }
        this.firstLimit = 100;
    }

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
            fetch('http://rt-dev.kbb1.com:8080/admin/rest/files' +
                  '?offset=' + startIndex +
                      '&limit=' + limit +
                      '&query=' + searchText)
                  .then((response) => {
                      if (!response.ok) {
                          throw Error('Error loading files, response not ok.');
                      }
                      this.setState({loadingFiles: false});
                      return response.json().then(json => {
                          if (json.status && json.status === 'ok') {
                              this.setState((prevState) => {
                                  const newFiles = prevState.resetFiles ? [] : prevState.files;
                                  json.files.forEach((f, i) => {
                                      newFiles[i + startIndex] = f;
                                      f.index = i + startIndex;
                                  });
                                  return {
                                      resetFiles: false,
                                      files: newFiles,
                                      matching: json.matching,
                                      total: json.total,
                                      error: '',
                                  };
                              });
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
        });
    }


    render() {
        const files = this.state.files;
        const rowGetter = ({ index }) => {
            return index < files.length && !!files[index] && !!files[index].id ? files[index] : {};
        }
        const isRowLoaded = ({ index }) => {
            return index < files.length && !!files[index] && files[index].id !== undefined;
        };
        const loadMoreRows = ({ startIndex, stopIndex }) => {
            return this.searchFiles(this.state.searchText, startIndex, stopIndex);
        }

        return (
            <div style={{ display: 'flex', flex: '1 1 auto', flexDirection: 'column'}}>
                <Header showRemoveIcon={this.state.showRemoveIcon}
                        searchText={this.state.searchText}
                        handleSearchChange={this.handleSearchChange}
                        handleSearchCancel={this.handleSearchCancel}
                        loadingFiles={this.state.loadingFiles}
                        error={this.state.error}
                        matching={this.state.matching}
                        total={this.state.total} />
                <div style={{ display: 'flex', flex: '1 1 auto', flexDirection: 'column'}}>
                    <InfiniteLoader ref="inf"
                                    isRowLoaded={isRowLoaded}
                                    threshold={100}
                                    loadMoreRows={loadMoreRows}
                                    rowCount={this.state.matching}>
                        {({ onRowsRendered, registerChild }) => (
                            <AutoSizer>
                                {({ width, height }) => (
                                    <Table headerHeight={50}
                                           height={height}
                                           width={width}
                                           rowCount={this.state.matching}
                                           ref={registerChild}
                                           onRowsRendered={onRowsRendered}
                                           rowRenderer={RowRenderer}
                                           rowGetter={rowGetter}
                                           rowHeight={50}>
                                        <Column label='Index'
                                                cellDataGetter={({ rowData }) => rowData.index}
                                                dataKey='index'
                                                width={60} />
                                        <Column label='ID'
                                                dataKey='id'
                                                cellRenderer={LinkToFileCellRenderer}
                                                width={80} />
                                        <Column label='UID'
                                                dataKey='uid'
                                                width={80} />
                                        <Column label='Name'
                                                dataKey='name'
                                                width={160} flexGrow={1}/>
                                        <Column label='Created at'
                                                dataKey='file_created_at'
                                                width={80}
                                                flexGrow={1} />
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
        <div className='ui fluid search'
             style={{ display: 'flex',
                      flexDirection: 'row',
                      justifyContent: 'space-between',
                      paddingLeft: 10,
                      paddingRight: 10 }}>
            <div>
                <div className='ui icon input'>
                    <input className='prompt'
                           type='text'
                           placeholder='Search files...'
                           value={props.searchText}
                           onChange={props.handleSearchChange} />
                    <i className='search icon' />
                </div>
                <i className='remove icon'
                   onClick={props.handleSearchCancel}
                   style={removeIconStyle} />
            </div>
            <div className='flex-space-between-center'>
                {props.loadingFiles &&
                    <span className='flex-space-between-center'>
                        <Spinner/>
                        <span style={{marginLeft: '10px'}}>Searching...</span>
                    </span>}
                {!!props.error && <span style={{color: 'red', marginLeft: '10px'}}>{props.error}</span>}
            </div>
            <div className='flex-space-between-center'>
                {props.matching >= 0 && props.total >= 0 &&
                 <span>Matched {props.matching} of {props.total}</span>}
            </div>
        </div>
    );
}

export default Files;
