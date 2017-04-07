import React from 'react';
import './Spinner.css'

const Spinner = (props) => {
    const size = props.size,
        style = {
            display: "inline-block",
            width: `${size}px`,
            height: `${size}px`,
            borderWidth: `${Math.max(Math.round(size / 10), 2)}px`
        };

    return <div className="spinner" style={style}/>;
};

Spinner.propTypes = {
    size: React.PropTypes.number
};

Spinner.defaultProps = {
    size: 24,
};

export default Spinner;
