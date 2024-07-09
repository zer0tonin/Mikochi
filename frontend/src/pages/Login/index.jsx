import { h } from "preact";
import { useContext, useEffect, useState } from "preact/hooks";
import "./style.css";

import { AuthContext } from "../../jwt";
import handleError from "../../error";

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
      const response = await fetch("/api/login", {
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
      <form onSubmit={onSubmit} class="form">
        <div class="container">
          <input
            type="text"
            value={username}
            onInput={(e) => setUsername(e.target.value)}
            class="input"
            id="username"
            placeholder="Username"
            aria-label="Username"
          />
          <i class="userIcon gg-user" />
        </div>
        <div class="container">
          <input
            type="password"
            value={password}
            onInput={(e) => setPassword(e.target.value)}
            class="input"
            id="password"
            placeholder="Password"
            aria-label="Password"
          />
          <i class="keyIcon gg-key" />
        </div>
        <div class="buttonContainer">
          <button type="submit" class="submit">
            Login
          </button>
        </div>
        {error !== "" && <div class="error">{error}</div>}
      </form>
    </main>
  );
};

export default Login;
