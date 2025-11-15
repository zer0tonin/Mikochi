import { h } from "preact";
import Download from "./download";
import { renameFilePath } from "./rename";
import { deleteFilePath } from "./delete";
import CopyLink from "./copylink";
import { useEffect, useRef, useState } from "preact/hooks";
import Icon from "../icon";

// DropDownMenu is slightly duplicated from the one found in menus because we need a different CSS
const DropDownMenu = ({ children }) => {
  const [visible, setVisible] = useState(false);
  const menuRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (menuRef.current && !menuRef.current.contains(event.target)) {
        setVisible(false);
      }
    };

    if (visible) {
      document.addEventListener("mousedown", handleClickOutside);
    } else {
      document.removeEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [visible]);

  return (
    <span
      ref={menuRef}
      class="actions-menu"
      onClick={() => setVisible(!visible)}
    >
      <Icon name="more-vertical-alt" title="actions" />
      {visible && <ul class="actions-dropdown-content">{children}</ul>}
    </span>
  );
};

export const DirActions = ({ filePath }) => {
  if (window.innerWidth < 768) {
    return (
      <DropDownMenu>
        <Download filePath={`${filePath}/`}>
          <li>Download</li>
        </Download>
        <li onClick={() => (renameFilePath.value = filePath)}>Rename</li>
        <li onClick={() => (deleteFilePath.value = filePath)}>Delete</li>
      </DropDownMenu>
    );
  }

  return (
    <>
      <Download filePath={`${filePath}/`}>
        <Icon name="arrow-down-o" title="Download" />
      </Download>
      <Icon
        name="rename"
        onClick={() => (renameFilePath.value = filePath)}
        title="Rename"
      />
      <Icon
        name="remove"
        onClick={() => (deleteFilePath.value = filePath)}
        title="Delete"
      />
    </>
  );
};

export const FileActions = ({ filePath }) => {
  if (window.innerWidth < 768) {
    return (
      <DropDownMenu>
        <Download filePath={filePath}>
          <li>Download</li>
        </Download>
        <CopyLink filePath={filePath}>
          <li>Copy stream link</li>
        </CopyLink>
        <li onClick={() => (renameFilePath.value = filePath)}>Rename</li>
        <li onClick={() => (deleteFilePath.value = filePath)}>Delete</li>
      </DropDownMenu>
    );
  }

  return (
    <>
      <Download filePath={filePath}>
        <Icon name="arrow-down-o" title="Download" />
      </Download>
      <CopyLink filePath={filePath}>
        <Icon name="copy" title="Copy stream link to clipboard" />
      </CopyLink>
      <Icon
        name="rename"
        onClick={() => (renameFilePath.value = filePath)}
        title="Rename"
      />
      <Icon
        name="remove"
        onClick={() => (deleteFilePath.value = filePath)}
        title="Delete"
      />
    </>
  );
};
