import React, { Component } from 'react';

class File extends Component {
    constructor(props) {
        super(props);

        this.state = {
            id: props.match.params.id,
        };
    }

    render() {
        return (
            <div>This is File {this.state.id}</div>
        );
    }

}

export default File;
