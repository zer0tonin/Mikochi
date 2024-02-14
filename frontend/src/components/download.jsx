import { h } from "preact";
import { useContext } from "preact/hooks";

import { AuthContext } from "../jwt";
import Icon from "./icon";

const downloadFile = (request, fileName) => {
  const link = document.createElement("a");
  link.href = request;
  link.download = fileName;
  link.click();
};

const Download = ({ filePath }) => {
  const { jwt } = useContext(AuthContext);

  const downloadWithAuth = async () => {
    const encodedFilePath = encodeURI(
      filePath.startsWith('/') ? filePath : `/${filePath}`
    );
    const response = await fetch(`/api/single-use?target=${encodedFilePath}`, {
      headers: {
        Accept: "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    });
    const json = await response.json();

    const auth = new URLSearchParams();
    auth.append("auth", json["token"]);
    downloadFile(
      `${window.location.protocol}//${window.location.hostname}${
        window.location.port == "" ? "" : `:${window.location.port}`
      }/api/stream${encodedFilePath}?${auth.toString()}`,
      filePath.split("/").at(-1),
    );
  };

  return (
    <a href="#" onClick={downloadWithAuth} style={{ color: "#E6E1C5" }}>
      <Icon name="arrow-down-o" title="Download" />
    </a>
  );
};

export default Download;
