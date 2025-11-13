import { h } from "preact";
import "./style.css";
import { useState } from "preact/hooks";
import Icon from "../icon";
import LogOut from "./logout";
import {mkdirOpen} from "../add/mkdir";
import {uploadOpen} from "../add/upload";

const Menu = ({ onHomeClick }) => {
  const [visible, setVisible] = useState(false);
  return (
    <span class="menu" onClick={() => setVisible(!visible)}>
      <Icon name="menu" />
      {visible && (
        <ul class="dropdown-content">
          <li onClick={onHomeClick}>Home</li>
          <li onClick={() => mkdirOpen.value = true}>New directory</li>
          <li onClick={() => uploadOpen.value = true}>Upload</li>
          <LogOut />
        </ul>
      )}
    </span>
  );
};

export default Menu;
