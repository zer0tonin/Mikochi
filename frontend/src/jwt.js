import { createContext } from "preact";

export const AuthContext = createContext(null);

// refreshJWT will get the jwt stored in localStorage and use it to obtain a fresh token
export const refreshJWT = (setJWT) => {
  const refresh = async () => {
    const storedToken = window.localStorage.getItem("jwt");

    // we always check refresh, to handle NO_AUTH=true backends

    const response = await fetch(`/api/refresh`, {
      headers: {
        Accept: "application/json",
        Authorization: `Bearer ${storedToken}`,
      },
    });

    if (response.status !== 200) {
      // we most likely have an expired token
      window.localStorage.removeItem("jwt");
      setJWT(null);
      return;
    }

    const json = await response.json();
    window.localStorage.setItem("jwt", json["token"]);
    setJWT(json["token"]);
  };

  refresh();
};
