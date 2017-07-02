import React from 'react'
import { Router, Route, IndexRoute, hashHistory } from 'react-router'
import Application from './application.jsx'
import ga from './library/analytics/ga'

const trackPage = function () {
  const { location } = this.state
  const path = window.location.pathname + '#' + location.pathname + location.search
  ga.pageview(path)
}

const router = (
  <Router history={hashHistory} onUpdate={trackPage}>
    <Route path="/" component={ Application }>
      <IndexRoute getComponent={(location, callback) => {
        require.ensure([], (require) => {
          callback(null, require('./view/home/app.jsx').default)
        })
      }} />
    </Route>
  </Router>
)

export default router
