import { h } from "preact";
import { useContext } from "preact/hooks";

import { AuthContext } from "../utils/jwt";

import Icon from "./icon";

// copyToClipboard will copy text to the clipboard using navigator.clipboard if available
// or fallback to document.execCommand
const copyToClipboard = async (textToCopy) => {
  // Navigator clipboard api needs a secure context (https)
  if (navigator.clipboard && window.isSecureContext) {
    await navigator.clipboard.writeText(textToCopy);
  } else {
    // Use the 'out of viewport hidden text area' trick
    const textArea = document.createElement("textarea");
    textArea.value = textToCopy;

    // Move textarea out of the viewport so it's not visible
    textArea.style.position = "absolute";
    textArea.style.left = "-999999px";

    document.body.prepend(textArea);
    textArea.select();

    try {
      document.execCommand("copy");
    } catch (error) {
      console.error(error);
    } finally {
      textArea.remove();
    }
  }
};

const CopyLink = ({ filePath }) => {
  const { jwt } = useContext(AuthContext);

  const copyWithAuth = async () => {
    const target = new URLSearchParams();
    target.append("target", filePath);
    const response = await fetch(`/api/single-use?${target.toString()}`, {
      headers: {
        Accept: "application/json",
        Authorization: `Bearer ${jwt}`,
      },
    });
    const json = await response.json();

    const auth = new URLSearchParams();
    auth.append("auth", json["token"]);
    await copyToClipboard(
      `${window.location.protocol}//${window.location.hostname}${
        window.location.port == "" ? "" : `:${window.location.port}`
      }/api/stream${filePath}?${auth.toString()}`
    );
  };
  return <Icon name="copy" onClick={copyWithAuth} title="Copy stream link" />;
};

export default CopyLink;
