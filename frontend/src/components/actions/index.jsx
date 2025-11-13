import { h } from "preact";
import Download from "./download";
import { RenameIcon, RenameListItem } from "./rename";
import { DeleteIcon, DeleteListItem } from "./delete";
import CopyLink from "./copylink";
import { useEffect, useRef, useState } from "preact/hooks";
import Icon from "../icon";

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

export const DirActions = ({ filePath, refresh, setRefresh }) => {
  if (window.innerWidth < 768) {
    return (
      <DropDownMenu>
        <Download filePath={`${filePath}/`}>
          <li>Download</li>
        </Download>
        <RenameListItem
          filePath={filePath}
          refresh={refresh}
          setRefresh={setRefresh}
        />
        <DeleteListItem
          filePath={filePath}
          refresh={refresh}
          setRefresh={setRefresh}
        />
      </DropDownMenu>
    );
  }

  return (
    <>
      <Download filePath={`${filePath}/`}>
        <Icon name="arrow-down-o" title="Download" />
      </Download>
      <RenameIcon
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
      />
      <DeleteIcon
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
      />
    </>
  );
};

export const FileActions = ({ filePath, refresh, setRefresh }) => {
  return (
    <>
      <Download filePath={`${filePath}/`}>
        <Icon name="arrow-down-o" title="Download" />
      </Download>
      <CopyLink filePath={filePath} />
      <RenameIcon
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
      />
      <DeleteIcon
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
      />
    </>
  );
};
