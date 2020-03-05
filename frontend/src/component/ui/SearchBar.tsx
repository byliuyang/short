import React, { Component, ChangeEvent } from 'react';
import classNames from 'classnames';

import './SearchBar.scss';
import { Url } from '../../entity/Url';
import { DebounceInput } from 'react-debounce-input';

interface State {
  showSuggestionBox: boolean;
}

interface Props {
  onChange: (text: String) => void;
  autoCompleteSuggestions?: Array<Url>;
}

export class SearchBar extends Component<Props, State> {
  state = {
    showSuggestionBox: true
  };

  handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange(event.target.value);
  };

  createAutoCompleteBox() {
    if (!this.props.autoCompleteSuggestions) {
      return <div />;
    }

    return (
      <ul
        className={classNames('suggestions', {
          show: this.state.showSuggestionBox
        })}
      >
        {this.props.autoCompleteSuggestions.map(e => (
          <li key={e.alias}>
            <a href={e.originalUrl}>{e.alias}</a>
          </li>
        ))}
      </ul>
    );
  }

  hideSuggestionsBox = () => {
    this.setState({
      showSuggestionBox: false
    });
  };

  showSuggestionsBox = () => {
    this.setState({
      showSuggestionBox: true
    });
  };

  render() {
    return (
      <div className="search-box">
        <div className="search-input">
          {/* Remove dependency on react-debounce-input
            TODO(issue#520): [Refactor] Implement debouncing for input in search bar */}
          <DebounceInput
            minLength={2}
            maxLength={50}
            placeholder={'Search short links'}
            debounceTimeout={300}
            onChange={this.handleChange}
            onFocus={this.showSuggestionsBox}
            onBlur={this.hideSuggestionsBox}
          />
          <i className="material-icons search">search</i>
        </div>
        {this.createAutoCompleteBox()}
      </div>
    );
  }
}