import { h } from "preact";
import { useContext } from "preact/hooks";

import { AuthContext } from "../utils/jwt";

import Icon from "./icon";

const CopyLink = ({ filePath }) => {
  const { jwt } = useContext(AuthContext);

  const copyToClipboard = async () => {
    const response = await fetch(`/api/single-use`, {
      headers: {
        Accept: "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    });
    const json = await response.json();

    const params = new URLSearchParams();
    params.append("auth", json["token"]);

    navigator.clipboard.writeText(
      `${window.location.protocol}//${window.location.hostname}${
        window.location.port == "" ? "" : `:${window.location.port}`
      }/api/stream${filePath}?${params.toString()}`
    );
  };
  return <Icon name="copy" onClick={copyToClipboard} />;
};

export default CopyLink;
