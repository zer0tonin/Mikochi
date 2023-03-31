import { h } from "preact";
import { Router } from "preact-router";

// Code-splitting is automated for `routes` directory
import Directory from "../routes/directory/directory";

const App = () => (
  <div id="app">
    <Router>
      <Directory path="/" />
      <Directory path="/:dirPath+" />
    </Router>
  </div>
);

export default App;
