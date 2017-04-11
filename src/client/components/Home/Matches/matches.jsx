import React, { Component } from 'react';
import { connect } from 'react-redux';
import * as actions from '../../../actions/matches.jsx';
import Match from './match.jsx';
import { axiosInstance } from '../../../actions/index.jsx';

class Matches extends Component {
  constructor(props) {
    super(props);
    this.state = {
      page: 0
    }
    this.previous = this.previous.bind(this);
    this.next = this.next.bind(this);
    this.fetchMatches = this.fetchMatches.bind(this);
    this.handleMatch = this.handleMatch.bind(this);
  }

  componentDidMount() {
    this.fetchMatches();
  }

  fetchMatches() {
    this.props.getMatches(this.props.email, this.props.city);
  }

  handleMatch(match, target, path) {
    let obj = {
      User: {
        Email: this.props.email
      }
    };
    obj[target] = {
      Email: match.Email
    };
    axiosInstance.post(path, obj)
                 .then(res => {
                   this.props.getMatches(this.props.email, this.props.city);
                 })
                 .catch(err => {
                   console.error(`cound not ${target} friend`);
                 });
  }

  addFriend(friend) {
    this.handleMatch(friend, "Friend", '/api/friends');
  }

  rejectMatch(reject) {
    this.handleMatch(reject, "Reject", '/api/reject');
  }

  saveMatch(match) {
    this.handleMatch(match, "Save", '/api/save');
  }

  previous() {
    this.setState({
      page: this.state.page - 1
    });
  }

  next() {
    this.setState({
      page: this.state.page + 1
    });
  }

  render() {
    return (
      <div>
        <div className="matches">
          { this.props.matches ? this.props.matches.slice(this.state.page * 4, this.state.page*4 + 4).map((match, i) => <Match match={ match } key={i} addFriend={this.addFriend.bind(this, match)} rejectMatch={this.rejectMatch.bind(this,match)} saveMatch={this.saveMatch.bind(this,match)}/>) : <div></div> }
          <div className="home_buttons">
            {this.state.page ? <input type= "image" onClick={this.previous} src={"./styles/left-arrow.svg"} className="arrow" /> : null }
            {this.props.matches.slice(this.state.page * 4, this.state.page*4 + 4).length >= 4 ? <input type= "image" onClick={this.next} src={"./styles/right-arrow.svg"} className="arrow" /> : null}
         </div>
        </div>
      </div>
    )
  }
}

function mapStateToProps(state){
  return {
    email: state.userInfo.user.Email,
    city: state.userInfo.user.City,
    matches: state.matches
  }
}

export default connect(mapStateToProps, actions)(Matches);