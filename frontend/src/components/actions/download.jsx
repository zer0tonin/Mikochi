import { h } from "preact";
import { useContext } from "preact/hooks";
import { AuthContext } from "../../jwt";
import {getApiUrl, getStreamUrl} from "../../url";

const downloadFile = (request, fileName) => {
  const link = document.createElement("a");
  link.href = request;
  link.download = fileName;
  link.click();
};

const Download = ({ filePath, children }) => {
  const { jwt } = useContext(AuthContext);

  const downloadWithAuth = async () => {
    const target = new URLSearchParams();
    target.append("target", filePath);
    const response = await fetch(getApiUrl(`stream-token?${target.toString()}`), {
      headers: {
        Accept: "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    });
    const json = await response.json();

    const auth = new URLSearchParams();
    auth.append("auth", json["token"]);
    downloadFile(
      getStreamUrl(filePath, auth),
      filePath.split("/").at(-1),
    );
  };

  return (
    <a
      href="#"
      onClick={downloadWithAuth}
      style={{ color: "#E6E1C5", textDecoration: "none" }}
    >
      {children}
    </a>
  );
};

export default Download;
