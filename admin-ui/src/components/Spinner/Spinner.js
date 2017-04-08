import React from 'react';
import PropTypes from 'prop-types';
import './Spinner.css'

const Spinner = (props) => {
    const { size } = props;
    const style = {
        display: "inline-block",
        width: `${size}px`,
        height: `${size}px`,
        borderWidth: `${Math.max(Math.round(size / 10), 2)}px`
    };

    return <div className="spinner" style={style}/>;
};

Spinner.propTypes = {
    size: PropTypes.number
};

Spinner.defaultProps = {
    size: 24,
};

export default Spinner;
