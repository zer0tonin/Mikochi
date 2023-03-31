import { h } from "preact";
import { Router } from "preact-router";

// Code-splitting is automated for `routes` directory
import Directory from "../routes/directory/directory";

// TODO: clean up dead code & prettify
const App = () => (
  <div id="app">
    <Router>
      <Directory path="/" />
      <Directory path="/:dirPath+" />
    </Router>
  </div>
);

export default App;
