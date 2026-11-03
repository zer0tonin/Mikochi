import { h } from "preact";
import { useContext } from "preact/hooks";
import { AuthContext } from "../../jwt";
import { getApiUrl } from "../../url";
import { addPrefix } from "../../location";

const LogOut = () => {
  const { jwt } = useContext(AuthContext);
  const onLogOut = async () => {
    try {
      if (!jwt) {
        console.error("No token found in context");
        return;
      }
      await fetch(getApiUrl("logout"), {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${jwt}`,
        },
      });

      window.location.href = addPrefix("/");
    } catch (error) {
      console.error("Error logging out:", error);
    } finally {
      localStorage.removeItem("jwt");
    }
  };

  return <li onClick={onLogOut}>Log out</li>;
};

export default LogOut;
