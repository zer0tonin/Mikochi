import { h } from 'preact';

const iconStyle = {
	margin: "0.1em",
	display: "inline-block",
	cursor: "pointer",
};

const Icon = ({ name, onClick }) => (
	<i style={iconStyle} class={`gg-${name}`} onClick={onClick}></i>
);

export default Icon;
