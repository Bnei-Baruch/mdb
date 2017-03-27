import React, { Component } from 'react';

class Welcome extends Component {
    render() {
        return (
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                margin: '20px',
            }}>
            <h2>Welcome to MDB Admin!</h2>
            </div>
        );
    }

}

export default Welcome;
