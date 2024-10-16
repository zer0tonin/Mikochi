import { h } from "preact";

import { useLocation } from "preact-iso";
import "./style.css";

const Logout = () => {
  const location = useLocation();
  const handleLogout = async () => {
    try {
      const token = localStorage.getItem("jwt");
      if (!token) {
        console.error("No token found in local storage");
        return;
      }
      await fetch("/api/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });

      localStorage.removeItem("jwt");

      //   location.route("/");
      window.location.href = "/";
    } catch (error) {
      console.error("Error logging out:", error);
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
