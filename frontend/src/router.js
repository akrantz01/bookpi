import React, { Component } from 'react'
import { HashRouter, Switch, Route } from 'react-router-dom'
import { ToastContainer, toast, Flip } from 'react-toastify'

import { Users } from './api'

import Header from './components/header'
import Footer from './components/footer'

import Home from './pages/home'
import SignIn from './pages/sign-in'
import SignUp from './pages/sign-up'
import NotFound from './pages/not-found'
import Chat from './pages/chat'
import Files from './pages/files'
import Shared from './pages/shared'
import Account from './pages/account'

class Router extends Component {
  constructor (props) {
    super(props)

    this.state = {
      loggedIn: false,
      loading: true,
      user: {}
    }
  }

  componentDidMount () {
    Users.readSelf().then(data => {
      if (data.status === 200) this.setState({ loggedIn: true, user: data.data })
      else if (data.status !== 401) toast.error(`Failed to read user data: (${data.status}) ${data.reason}`)
    })
      .finally(() => this.setState({ loading: false }))
  }

    updateUser = name => this.setState({ user: { ...this.state.user, name } });

    login () {
      this.setState({ loggedIn: true, loading: true })
      Users.readSelf().then(data => {
        if (data.status === 200) this.setState({ user: data.data, loading: false })
        else if (data.status !== 401) toast.error(`Failed to read user data: (${data.status}) ${data.reason}`)
      })
    }

    logout = () => this.setState({ loggedIn: false, user: {} });

    render () {
      if (this.state.loading) {
        return (
          <div className="spinner-border text-info" style={{ position: 'absolute', top: '40%', right: '50%' }} role="status">
            <span className="sr-only">Loading...</span>
          </div>
        )
      }

      return (
        <HashRouter>
          <ToastContainer position="bottom-right" autoClose={4000} closeOnClick pauseOnHover draggable transition={Flip}/>
          <Header loggedIn={this.state.loggedIn} logout={this.logout.bind(this)}/>
          <main role="main" className="flex-shrink-0" style={{ marginTop: '40px' }}>
            <Switch>
              <Route path="/" exact><Home/></Route>
              <Route path="/sign-in" exact><SignIn login={this.login.bind(this)} loggedIn={this.state.loggedIn} /></Route>
              <Route path="/sign-up" exact><SignUp loggedIn={this.state.loggedIn} /></Route>
              <Route path="/chat" exact><Chat loggedIn={this.state.loggedIn} username={this.state.user.username} /></Route>
              <Route path="/my-files" exact><Files loggedIn={this.state.loggedIn} /></Route>
              <Route path="/shared-with-me" exact><Shared loggedIn={this.state.loggedIn} /></Route>
              <Route path="/account" exact><Account loggedIn={this.state.loggedIn} user={this.state.user} updateName={this.updateUser.bind(this)} /></Route>
              <Route path="*"><NotFound/></Route>
            </Switch>
          </main>
          <Footer/>
        </HashRouter>
      )
    }
}

export default Router
