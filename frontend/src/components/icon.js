import { h } from "preact";

const iconStyle = {
  display: "inline-block",
  cursor: "pointer",
  marginLeft: "0.4em",
};

const Icon = ({ name, onClick }) => (
  <i style={iconStyle} class={`gg-${name}`} onClick={onClick} />
);

export default Icon;
