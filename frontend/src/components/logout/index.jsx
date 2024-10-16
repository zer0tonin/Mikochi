import { h } from "preact";
import React, { useContext } from "react";
import { AuthContext } from "../../jwt";

import { useLocation } from "preact-iso";
import "./style.css";

const Logout = () => {
  const location = useLocation();
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

      location.route("/");
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
