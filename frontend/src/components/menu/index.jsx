import { h } from "preact";
import "./style.css";
import { useEffect, useRef, useState } from "preact/hooks";
import Icon from "../icon";
import LogOut from "./logout";
import { mkdirOpen } from "../add/mkdir";
import { uploadOpen } from "../add/upload";

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
    <span ref={menuRef} class="menu" onClick={() => setVisible(!visible)}>
      <Icon name={"menu"} />
      {visible && <ul class="dropdown-content">{children}</ul>}
    </span>
  );
};

const Menu = ({ onHomeClick }) => {
  return (
    <DropDownMenu>
      <li onClick={onHomeClick}>Home</li>
      <li onClick={() => (mkdirOpen.value = true)}>New directory</li>
      <li onClick={() => (uploadOpen.value = true)}>Upload</li>
      <LogOut />
    </DropDownMenu>
  );
};

export default Menu;
