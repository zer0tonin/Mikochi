import { h } from "preact";
import "./icon.css";

const iconStyle = {
  display: "inline-block",
  cursor: "pointer",
  marginLeft: "0.4em",
};

const Icon = ({ name, title, onClick }) => (
  <i
    style={iconStyle}
    class={`gg-${name}`}
    onClick={onClick}
    title={title}
    aria-label={title}
  />
);

const bigIconStyle = {
  display: "inline-block",
  cursor: "pointer",
  "--ggs": 3,
};

export const BigIcon = ({ name, title, onClick }) => (
  <i
    style={bigIconStyle}
    class={`gg-${name}`}
    onClick={onClick}
    title={title}
    aria-label={title}
  />
);

export default Icon;
