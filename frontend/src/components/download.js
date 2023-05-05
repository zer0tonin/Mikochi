import { h } from "preact";

import Icon from "./icon";

const Download = ({ filePath }) => {
  return (
    <a
      href={`${window.location.protocol}//${window.location.hostname}${
        window.location.port == "" ? "" : `:${window.location.port}`
      }/api/stream${filePath}`}
      download={filePath.split("/").pop()}
      style={{ color: "#E6E1C5" }}
    >
      <Icon name="arrow-down-o" />
    </a>
  );
};

export default Download;
