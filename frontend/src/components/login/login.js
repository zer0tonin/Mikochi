import { h } from "preact";
import { useContext, useState } from "preact/hooks";
import style from "./style.css";

import { AuthContext } from "../../utils/jwt";

const Login = () => {
  const { setJWT } = useContext(AuthContext);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

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
        setError(json["err"]);
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
      <form onSubmit={onSubmit} class={style.form}>
        <div class={style.container}>
          <input
            type="text"
            value={username}
            onInput={(e) => setUsername(e.target.value)}
            class={style.input}
            id="username"
            placeholder="Username"
          />
          <i class={`${style.userIcon} gg-user`} />
        </div>
        <div class={style.container}>
          <input
            type="password"
            value={password}
            onInput={(e) => setPassword(e.target.value)}
            class={style.input}
            id="password"
            placeholder="Password"
          />
          <i class={`${style.keyIcon} gg-key`} />
        </div>
        <div class={style.buttonContainer}>
          <button type="submit" class={style.submit}>
            Login
          </button>
        </div>
        {error !== "" && <div class={style.error}>{error}</div>}
      </form>
    </main>
  );
};

export default Login;
