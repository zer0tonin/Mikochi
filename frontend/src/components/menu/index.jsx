import { h } from "preact";
import "./style.css";
import {useState} from "preact/hooks";
import Icon from "../icon";
import LogOut from "./logout";


const Menu = ({ onHomeClick }) => {
  const [visible, setVisible] = useState(false)
  return (
    <span class="menu" onClick={() => setVisible(!visible)}>
      <Icon name="menu"  />
      {visible &&
        <ul class="dropdown-content">
          <li onClick={onHomeClick}>Home</li>
          <LogOut />
        </ul>}
    </span>
  );
}

export default Menu;
