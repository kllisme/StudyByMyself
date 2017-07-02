import React from "react"
import "./application.less"

const Application = React.createClass({
 render() {
    const { children } = this.props
    return (<div className="application">
      <div className="application-main">
        {children}
      </div>
    </div>)
  }
})

export default Application
