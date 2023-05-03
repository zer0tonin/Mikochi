import { h, createContext } from "preact";
import { useState, useEffect, useMemo } from "preact/hooks";
import { Router } from "preact-router";

import { refreshJWT } from "../utils/jwt"

// Code-splitting is automated for `routes` directory
import Directory from "../routes/directory/directory";
import Login from "../routes/login/login";

const AuthContext = createContext()

const App = () => {
  const [jwt, setJWT] = useState(null)
  useEffect(() => refreshJWT(setJWT), [])
  const auth = useMemo(() => {
    return {jwt, setJWT}
  }, [jwt])

  return (
    <AuthContext.Provider value={auth}>
      <div id="app">
        { jwt === null
          ? <Login />
          : (
              <Router>
                <Directory path="/" />
                <Directory path="/:dirPath+" />
              </Router>
        )}
      </div>
    </AuthContext.Provider>
  );
}

export default App;
