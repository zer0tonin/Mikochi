import { h } from "preact";

const iconStyle = {
  display: "inline-block",
  cursor: "pointer",
  marginLeft: "0.4em",
};

const Icon = ({ name, title, onClick }) => (
  <i style={iconStyle} class={`gg-${name}`} onClick={onClick} title={title} aria-label={title} />
);

export default Icon;
