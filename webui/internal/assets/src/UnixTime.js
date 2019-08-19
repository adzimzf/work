import React from 'react';
import PropTypes from 'prop-types';

export default class UnixTime extends React.Component {
  static propTypes = {
    ts: PropTypes.number.isRequired,
  };


  render() {
    let t = new Date(this.props.ts * 1e3);
    return (
      <time dateTime={t.toUTCString()}>{t.toString()}</time>
    );
  }
}
