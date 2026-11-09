import { h } from "preact";
import { useContext, useEffect, useState } from "preact/hooks";
import "./style.css";

import { AuthContext } from "../../jwt";
import handleError from "../../error";
import { getApiUrl } from "../../url";

const Login = () => {
  const { setJWT } = useContext(AuthContext);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    document.title = `Mikochi - Login`;
  }, []);

  const onSubmit = (e) => {
    e.preventDefault();
    const postLogin = async () => {
      const response = await fetch(getApiUrl("login"), {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });
      const json = await response.json();
      if (response.status !== 200) {
        setError(await handleError(response));
        return;
      }
      window.localStorage.setItem("jwt", json["token"]);
      setJWT(json["token"]);
      setError("");
    };
    postLogin();
  };

  return (
    <main>
      <form onSubmit={onSubmit} class="login-form">
        <div class="login-container">
          <input
            type="text"
            value={username}
            onInput={(e) => setUsername(e.target.value)}
            class="login-input"
            id="username"
            placeholder="Username"
            aria-label="Username"
          />
          {
            //<i class="login-userIcon gg-user" />
          }
        </div>
        <div class="login-container">
          <input
            type="password"
            value={password}
            onInput={(e) => setPassword(e.target.value)}
            class="login-input"
            id="password"
            placeholder="Password"
            aria-label="Password"
          />
          {
            //<i class="login-keyIcon gg-key" />
          }
        </div>
        <div class="login-buttonContainer">
          <button type="submit" class="login-submit">
            Login
          </button>
        </div>
        {error !== "" && <div class="login-error">{error}</div>}
      </form>
    </main>
  );
};

export default Login;
