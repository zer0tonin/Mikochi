import { h } from "preact";
import React, { useContext } from "react";
import { AuthContext } from "../../jwt";

import "./style.css";

const Logout = () => {
  const { jwt } = useContext(AuthContext);

  const handleLogout = async () => {
    try {
      if (!jwt) {
        console.error("No token found in context");
        return;
      }
      await fetch("/api/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${jwt}`,
        },
      });

      window.location.href = "/";
    } catch (error) {
      console.error("Error logging out:", error);
    } finally {
      localStorage.removeItem("jwt");
    }
  };

  return (
    <>
      <button class="logout" onClick={handleLogout}>
        Log out
      </button>
    </>
  );
};

export default Logout;
